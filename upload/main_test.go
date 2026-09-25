package main

import "testing"

func TestValidateManuscript_MarkdownFileIsAccepted(t *testing.T) {
	path := "manuscript.md"

	got, err := ValidateManuscript(path)
	if err != nil {
		t.Fatalf("ValidateManuscript(%q) returned unexpected error: %v", path, err)
	}

	if got != path {
		t.Errorf("ValidateManuscript(%q) = %q, want %q", path, got, path)
	}
}

func TestValidateManuscript_NonMarkdownFileIsRejected(t *testing.T) {
	path := "manuscript.docx"

	_, err := ValidateManuscript(path)
	if err == nil {
		t.Fatalf("ValidateManuscript(%q) expected an error, got none", path)
	}
}
