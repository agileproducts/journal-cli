package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateManuscript checks that the given file has a supported extension
// and returns the path unchanged so it can be piped to a successor program.
func ValidateManuscript(path string) (string, error) {
	if filepath.Ext(path) != ".md" {
		return "", fmt.Errorf("only markdown (.md) files are supported at present")
	}

	return path, nil
}

func main() {
	path, err := ResolvePath(os.Args, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "usage: jcli-upload <manuscript-path> (or pipe a path via stdin)")
		os.Exit(1)
	}

	path, err = ValidateManuscript(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(path)
}
