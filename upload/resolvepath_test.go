package main

import (
	"strings"
	"testing"
)

func TestResolvePath_FromArgs(t *testing.T) {
	got, err := ResolvePath([]string{"jcli-upload", "manuscript.md"}, strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "manuscript.md" {
		t.Errorf("ResolvePath() = %q, want %q", got, "manuscript.md")
	}
}

func TestResolvePath_FromStdinWhenNoArgs(t *testing.T) {
	got, err := ResolvePath([]string{"jcli-upload"}, strings.NewReader("manuscript.md\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "manuscript.md" {
		t.Errorf("ResolvePath() = %q, want %q", got, "manuscript.md")
	}
}

func TestResolvePath_ArgsTakePriorityOverStdin(t *testing.T) {
	got, err := ResolvePath([]string{"jcli-upload", "from-args.md"}, strings.NewReader("from-stdin.md\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "from-args.md" {
		t.Errorf("ResolvePath() = %q, want %q", got, "from-args.md")
	}
}

func TestResolvePath_ErrorsWhenNothingProvided(t *testing.T) {
	_, err := ResolvePath([]string{"jcli-upload"}, strings.NewReader(""))
	if err == nil {
		t.Fatal("expected an error when no path is provided via args or stdin")
	}
}
