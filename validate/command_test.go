package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jcli-submission"
)

func TestRunPersistsPassingAndFailingValidation(t *testing.T) {
	for _, tc := range []struct {
		name, content string
		wantCode      int
	}{
		{"passing", "# Title\n## Abstract\n## References\n", 0},
		{"failing", "## Abstract\n", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := submission.Store{Root: t.TempDir()}
			path := filepath.Join(t.TempDir(), "paper.md")
			if err := os.WriteFile(path, []byte(tc.content), 0600); err != nil {
				t.Fatal(err)
			}
			r, err := s.Create(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			for _, args := range [][]string{{"jcli-validate", r.ID}, {"jcli-validate"}} {
				var out, diagnostics bytes.Buffer
				input := r.ID + "\n"
				if len(args) > 1 {
					input = "ignored\n"
				}
				code := run(args, strings.NewReader(input), &out, &diagnostics, s)
				if code != tc.wantCode {
					t.Fatalf("exit %d: %s", code, &diagnostics)
				}
				wantOut := ""
				if tc.wantCode == 0 {
					wantOut = r.ID + "\n"
				}
				if out.String() != wantOut || !strings.Contains(diagnostics.String(), "✓ Abstract") {
					t.Fatalf("stdout=%q stderr=%q", &out, &diagnostics)
				}
				if tc.wantCode == 1 && (!strings.Contains(diagnostics.String(), "✗ Title") || !strings.Contains(diagnostics.String(), "✗ References")) {
					t.Fatalf("missing failed checks: %s", &diagnostics)
				}
			}
			loaded, err := s.Load(r.ID)
			if err != nil {
				t.Fatal(err)
			}
			if len(loaded.Validations) != 2 || len(loaded.Validations[0].Checks) != 3 || loaded.Validations[0].Checks[0].Passed != (tc.wantCode == 0) {
				t.Fatalf("unexpected history: %+v", loaded.Validations)
			}
		})
	}
}

func TestRunValidationInputErrors(t *testing.T) {
	for _, input := range []string{"", "paper.md\n", "sub-00000000000000000000000000000000\n"} {
		var out, diagnostics bytes.Buffer
		if code := run([]string{"jcli-validate"}, strings.NewReader(input), &out, &diagnostics, submission.Store{Root: t.TempDir()}); code != 1 || out.Len() != 0 || diagnostics.Len() == 0 {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, &out, &diagnostics)
		}
	}
}

func TestRunChecksAfterLongParagraph(t *testing.T) {
	content := "# Title\n" + strings.Repeat("a", 70000) + "\n## Abstract\n## References\n"
	if !AllPassed(RunChecks(content)) {
		t.Fatal("a long paragraph prevented later headings from being checked")
	}
}

func TestRunDoesNotReportSuccessWhenSavingFails(t *testing.T) {
	s := submission.Store{Root: t.TempDir()}
	path := filepath.Join(t.TempDir(), "paper.md")
	if err := os.WriteFile(path, []byte("# Title\n## Abstract\n## References\n"), 0600); err != nil {
		t.Fatal(err)
	}
	r, err := s.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	// A directory at the lock path makes saving fail without permission changes.
	if err := os.Mkdir(filepath.Join(s.Root, "submissions", r.ID, ".lock"), 0700); err != nil {
		t.Fatal(err)
	}
	var out, diagnostics bytes.Buffer
	if code := run([]string{"jcli-validate", r.ID}, strings.NewReader(""), &out, &diagnostics, s); code != 1 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "save validation") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, &out, &diagnostics)
	}
}
