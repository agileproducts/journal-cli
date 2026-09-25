package main

import "testing"

func findCheck(checks []Check, name string) (Check, bool) {
	for _, c := range checks {
		if c.Name == name {
			return c, true
		}
	}
	return Check{}, false
}

func TestRunChecks_AllPresent(t *testing.T) {
	content := "# My Manuscript Title\n\nSome intro text.\n\n## Abstract\n\nThis is the abstract.\n\n## References\n\n1. Someone, 2024.\n"

	checks := RunChecks(content)

	for _, name := range []string{"Title", "Abstract", "References"} {
		c, ok := findCheck(checks, name)
		if !ok {
			t.Fatalf("expected a check named %q, got %v", name, checks)
		}
		if !c.Passed {
			t.Errorf("expected check %q to pass, but it failed", name)
		}
	}
}

func TestRunChecks_MissingTitle(t *testing.T) {
	content := "Some intro text without a heading.\n\n## Abstract\n\nThis is the abstract.\n\n## References\n\n1. Someone, 2024.\n"

	checks := RunChecks(content)

	c, ok := findCheck(checks, "Title")
	if !ok {
		t.Fatalf("expected a check named %q", "Title")
	}
	if c.Passed {
		t.Errorf("expected Title check to fail when there is no top-level heading")
	}
}

func TestRunChecks_MissingAbstract(t *testing.T) {
	content := "# My Manuscript Title\n\nSome intro text.\n\n## References\n\n1. Someone, 2024.\n"

	checks := RunChecks(content)

	c, ok := findCheck(checks, "Abstract")
	if !ok {
		t.Fatalf("expected a check named %q", "Abstract")
	}
	if c.Passed {
		t.Errorf("expected Abstract check to fail when there is no ## Abstract section")
	}
}

func TestRunChecks_MissingReferences(t *testing.T) {
	content := "# My Manuscript Title\n\nSome intro text.\n\n## Abstract\n\nThis is the abstract.\n"

	checks := RunChecks(content)

	c, ok := findCheck(checks, "References")
	if !ok {
		t.Fatalf("expected a check named %q", "References")
	}
	if c.Passed {
		t.Errorf("expected References check to fail when there is no ## References section")
	}
}

func TestAllPassed(t *testing.T) {
	passing := []Check{{Name: "Title", Passed: true}, {Name: "Abstract", Passed: true}}
	if !AllPassed(passing) {
		t.Errorf("expected AllPassed to be true when every check passed")
	}

	failing := []Check{{Name: "Title", Passed: true}, {Name: "Abstract", Passed: false}}
	if AllPassed(failing) {
		t.Errorf("expected AllPassed to be false when a check failed")
	}
}
