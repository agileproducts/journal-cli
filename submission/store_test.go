package submission

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func source(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "my manuscript.md")
	if err := os.WriteFile(p, []byte("# Original\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCreateSnapshotAndReload(t *testing.T) {
	s := Store{Root: t.TempDir()}
	p := source(t)
	r, err := s.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	loaded, data, err := (Store{Root: s.Root}).ReadManuscript(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# Original\n" || loaded.Manuscript.Revision != 1 || loaded.SchemaVersion != 1 || loaded.CreatedAt.IsZero() || loaded.Manuscript.OriginalName != "my manuscript.md" {
		t.Fatalf("unexpected snapshot: %+v, %q", loaded, data)
	}
	second, err := s.Create(p)
	if err != nil || second.ID == r.ID {
		t.Fatalf("IDs must be unique: %v", err)
	}
}

func TestCreateRejectsMissingAndNonRegularFiles(t *testing.T) {
	s := Store{Root: t.TempDir()}
	for _, p := range []string{filepath.Join(t.TempDir(), "missing.md"), t.TempDir()} {
		if _, err := s.Create(p); err == nil {
			t.Fatalf("accepted %q", p)
		}
	}
	entries, err := os.ReadDir(s.Root)
	if err != nil || len(entries) != 0 {
		t.Fatalf("failed upload left state: %v, %v", entries, err)
	}
}

func TestLoadRejectsInvalidMissingAndCorruptRecords(t *testing.T) {
	s := Store{Root: t.TempDir()}
	for _, id := range []string{"../escape", "/tmp/file", "", "sub-00000000000000000000000000000000"} {
		if _, err := s.Load(id); err == nil {
			t.Fatalf("accepted %q", id)
		}
	}
	r, err := s.Create(source(t))
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(s.Root, "submissions", r.ID, "submission.json")
	for _, content := range []string{"{broken", `{}`, `{"schema_version":999}`} {
		if err := os.WriteFile(p, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Load(r.ID); err == nil {
			t.Fatalf("accepted corrupt record %q", content)
		}
	}
}

func TestValidationHistoryAndStaleRevision(t *testing.T) {
	s := Store{Root: t.TempDir()}
	r, err := s.Create(source(t))
	if err != nil {
		t.Fatal(err)
	}
	checks := []Check{{Name: "Title", Passed: true}, {Name: "Abstract", Passed: false}}
	if err := s.RecordValidation(r.ID, 1, checks); err != nil {
		t.Fatal(err)
	}
	if err := s.RecordValidation(r.ID, 2, checks); err == nil {
		t.Fatal("accepted stale revision")
	}
	loaded, err := s.Load(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Validations) != 1 || loaded.Validations[0].Revision != 1 || loaded.Validations[0].CheckedAt.IsZero() || loaded.Validations[0].Checks[1].Passed {
		t.Fatalf("failure not saved: %+v", loaded)
	}
}

func TestConcurrentValidationUpdatesPreserveEveryAttempt(t *testing.T) {
	s := Store{Root: t.TempDir()}
	r, err := s.Create(source(t))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.RecordValidation(r.ID, 1, []Check{{Name: "Title", Passed: true}}); err != nil {
				t.Error(err)
			}
			if _, err := s.Load(r.ID); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	loaded, err := s.Load(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Validations) != 20 {
		t.Fatalf("lost updates: %d", len(loaded.Validations))
	}
}

func TestDefaultRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("JCLI_HOME", "")
	t.Setenv("XDG_DATA_HOME", "")
	root, err := DefaultRoot()
	if err != nil || root != filepath.Join(home, ".local", "share", "jcli") {
		t.Fatalf("root=%q: %v", root, err)
	}
	xdg := t.TempDir()
	t.Setenv("XDG_DATA_HOME", xdg)
	root, err = DefaultRoot()
	if err != nil || root != filepath.Join(xdg, "jcli") {
		t.Fatalf("root=%q: %v", root, err)
	}
	custom := t.TempDir()
	t.Setenv("JCLI_HOME", custom)
	root, err = DefaultRoot()
	if err != nil || root != custom {
		t.Fatalf("root=%q: %v", root, err)
	}
}
