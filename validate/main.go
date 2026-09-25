package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"jcli-resolvepath"
	"jcli-submission"
)

// Check represents the outcome of a single validation rule run against a
// manuscript.
type Check = submission.Check

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

	for line := range strings.Lines(content) {
		line = strings.TrimSpace(line)

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

func printChecks(stderr io.Writer, checks []Check) {
	color := false
	if f, ok := stderr.(*os.File); ok && os.Getenv("NO_COLOR") == "" {
		if info, err := f.Stat(); err == nil {
			color = info.Mode()&os.ModeCharDevice != 0
		}
	}
	for _, c := range checks {
		mark, tint := "✓", colorGreen
		if !c.Passed {
			mark, tint = "✗", colorRed
		}
		if color {
			mark = tint + mark + colorReset
		}
		fmt.Fprintf(stderr, "%s %s\n", mark, c.Name)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, store submission.Store) int {
	id, err := resolvepath.Resolve(args, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "usage: jcli-validate <submission-id> (or pipe an ID via stdin): %v\n", err)
		return 1
	}

	record, data, err := store.ReadManuscript(id)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	checks := RunChecks(string(data))
	if err := store.RecordValidation(id, record.Manuscript.Revision, checks); err != nil {
		fmt.Fprintf(stderr, "could not save validation: %v\n", err)
		return 1
	}
	printChecks(stderr, checks)

	if !AllPassed(checks) {
		fmt.Fprintf(stderr, "manuscript failed validation (results saved for %s)\n", id)
		return 1
	}
	if _, err := fmt.Fprintln(stdout, id); err != nil {
		fmt.Fprintf(stderr, "write submission ID: %v\n", err)
		return 1
	}
	return 0
}

func main() {
	root, err := submission.DefaultRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(run(os.Args, os.Stdin, os.Stdout, os.Stderr, submission.Store{Root: root}))
}
