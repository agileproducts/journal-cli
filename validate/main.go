package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"jcli-resolvepath"
)

// Check represents the outcome of a single validation rule run against a
// manuscript.
type Check struct {
	Name   string
	Passed bool
}

var (
	titleRe      = regexp.MustCompile(`^#\s+\S`)
	abstractRe   = regexp.MustCompile(`(?i)^##\s*abstract\b`)
	referencesRe = regexp.MustCompile(`(?i)^##\s*references\b`)
)

// RunChecks inspects manuscript content and reports whether it has a title,
// an abstract section, and a references section.
func RunChecks(content string) []Check {
	hasTitle := false
	hasAbstract := false
	hasReferences := false

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if titleRe.MatchString(line) {
			hasTitle = true
		}
		if abstractRe.MatchString(line) {
			hasAbstract = true
		}
		if referencesRe.MatchString(line) {
			hasReferences = true
		}
	}

	return []Check{
		{Name: "Title", Passed: hasTitle},
		{Name: "Abstract", Passed: hasAbstract},
		{Name: "References", Passed: hasReferences},
	}
}

// AllPassed reports whether every check in the slice passed.
func AllPassed(checks []Check) bool {
	for _, c := range checks {
		if !c.Passed {
			return false
		}
	}
	return true
}

const (
	colorReset = "\x1b[0m"
	colorGreen = "\x1b[32m"
	colorRed   = "\x1b[31m"
)

func printChecks(checks []Check) {
	for _, c := range checks {
		mark := colorGreen + "✓" + colorReset
		if !c.Passed {
			mark = colorRed + "✗" + colorReset
		}
		fmt.Fprintf(os.Stderr, "%s %s\n", mark, c.Name)
	}
}

func main() {
	path, err := resolvepath.Resolve(os.Args, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "usage: jcli-validate <manuscript-path> (or pipe a path via stdin)")
		os.Exit(1)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read %s: %v\n", path, err)
		os.Exit(1)
	}

	checks := RunChecks(string(data))
	printChecks(checks)

	if !AllPassed(checks) {
		fmt.Fprintln(os.Stderr, "manuscript failed validation")
		os.Exit(1)
	}

	fmt.Println(path)
}
