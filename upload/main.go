package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"jcli-resolvepath"
	"jcli-submission"
)

// ValidateManuscript checks that the given file has a supported extension
// and returns the source path to snapshot.
func ValidateManuscript(path string) (string, error) {
	if filepath.Ext(path) != ".md" {
		return "", fmt.Errorf("only markdown (.md) files are supported at present")
	}

	return path, nil
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, store submission.Store) int {
	path, err := resolvepath.Resolve(args, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "usage: jcli-upload <manuscript-path> (or pipe a path via stdin): %v\n", err)
		return 1
	}

	path, err = ValidateManuscript(path)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	record, err := store.Create(path)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if _, err := fmt.Fprintln(stdout, record.ID); err != nil {
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
