package resolvepath

import (
	"strings"
	"testing"
)

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
