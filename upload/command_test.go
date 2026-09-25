package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jcli-submission"
)

func TestRunCreatesSubmission(t *testing.T) {
	path := filepath.Join(t.TempDir(), "my manuscript.md")
	if err := os.WriteFile(path, []byte("# Title\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"jcli-upload", path}, {"jcli-upload"}} {
		s := submission.Store{Root: t.TempDir()}
		var out, diagnostics bytes.Buffer
		input := path + "\n"
		if len(args) > 1 {
			input = "ignored.md\n"
		}
		if code := run(args, strings.NewReader(input), &out, &diagnostics, s); code != 0 {
			t.Fatalf("exit %d: %s", code, &diagnostics)
		}
		id := strings.TrimSuffix(out.String(), "\n")
		_, data, err := s.ReadManuscript(id)
		if err != nil || string(data) != "# Title\n" {
			t.Fatalf("invalid upload stdout=%q: %v", out.String(), err)
		}
	}
}

func TestRunUploadErrorsHaveNoSuccessOutput(t *testing.T) {
	for _, args := range [][]string{{"jcli-upload"}, {"jcli-upload", "wrong.docx"}, {"jcli-upload", "missing.md"}, {"jcli-upload", t.TempDir()}} {
		var out, diagnostics bytes.Buffer
		if code := run(args, strings.NewReader(""), &out, &diagnostics, submission.Store{Root: t.TempDir()}); code != 1 || out.Len() != 0 || diagnostics.Len() == 0 {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, &out, &diagnostics)
		}
	}
}
