package resolvepath

import (
	"errors"
	"strings"
	"testing"
)

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("input failed") }

func TestResolveSubmissionIDAndFirstLine(t *testing.T) {
	got, err := Resolve([]string{"jcli-tool"}, strings.NewReader("sub-123\nignored\n"))
	if err != nil || got != "sub-123" {
		t.Fatalf("got %q: %v", got, err)
	}
}

func TestResolveEmptyArgumentDoesNotFallBackToStdin(t *testing.T) {
	if _, err := Resolve([]string{"jcli-tool", ""}, strings.NewReader("sub-123\n")); err == nil {
		t.Fatal("empty argument unexpectedly fell back to stdin")
	}
}

func TestResolveReportsInputErrors(t *testing.T) {
	_, err := Resolve(nil, failingReader{})
	if err == nil || !strings.Contains(err.Error(), "input failed") {
		t.Fatalf("lost reader error: %v", err)
	}
}

func TestResolve_FromArgs(t *testing.T) {
	got, err := Resolve([]string{"jcli-tool", "manuscript.md"}, strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "manuscript.md" {
		t.Errorf("Resolve() = %q, want %q", got, "manuscript.md")
	}
}

func TestResolve_FromStdinWhenNoArgs(t *testing.T) {
	got, err := Resolve([]string{"jcli-tool"}, strings.NewReader("manuscript.md\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "manuscript.md" {
		t.Errorf("Resolve() = %q, want %q", got, "manuscript.md")
	}
}

func TestResolve_ArgsTakePriorityOverStdin(t *testing.T) {
	got, err := Resolve([]string{"jcli-tool", "from-args.md"}, strings.NewReader("from-stdin.md\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "from-args.md" {
		t.Errorf("Resolve() = %q, want %q", got, "from-args.md")
	}
}

func TestResolve_ErrorsWhenNothingProvided(t *testing.T) {
	_, err := Resolve([]string{"jcli-tool"}, strings.NewReader(""))
	if err == nil {
		t.Fatal("expected an error when no path is provided via args or stdin")
	}
}
