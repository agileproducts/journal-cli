package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jcli-submission"
)

func TestRunInspectSavedRecordWithoutChangingIt(t *testing.T) {
	s := submission.Store{Root: t.TempDir()}
	path := filepath.Join(t.TempDir(), "paper.md")
	if err := os.WriteFile(path, []byte("# Title"), 0600); err != nil {
		t.Fatal(err)
	}
	r, err := s.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.RecordValidation(r.ID, 1, []submission.Check{{Name: "Abstract", Passed: false}}); err != nil {
		t.Fatal(err)
	}
	recordPath := filepath.Join(s.Root, "submissions", r.ID, "submission.json")
	before, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"jcli-inspect", r.ID}, {"jcli-inspect"}} {
		var out, diagnostics bytes.Buffer
		input := r.ID + "\n"
		if len(args) > 1 {
			input = "ignored\n"
		}
		if code := run(args, strings.NewReader(input), &out, &diagnostics, s); code != 0 {
			t.Fatalf("exit %d: %s", code, &diagnostics)
		}
		var got submission.Record
		if err := json.Unmarshal(out.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got.ID != r.ID || len(got.Validations) != 1 || got.Validations[0].Checks[0].Passed {
			t.Fatalf("unexpected record: %+v", got)
		}
	}
	after, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("inspect changed the record")
	}
}

func TestRunInspectErrors(t *testing.T) {
	for _, input := range []string{"", "invalid\n", "sub-00000000000000000000000000000000\n"} {
		var out, diagnostics bytes.Buffer
		if code := run([]string{"jcli-inspect"}, strings.NewReader(input), &out, &diagnostics, submission.Store{Root: t.TempDir()}); code != 1 || out.Len() != 0 || diagnostics.Len() == 0 {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, &out, &diagnostics)
		}
	}
}
