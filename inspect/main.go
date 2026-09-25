package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"jcli-resolvepath"
	"jcli-submission"
)

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, store submission.Store) int {
	id, err := resolvepath.Resolve(args, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "usage: jcli-inspect <submission-id> (or pipe an ID via stdin): %v\n", err)
		return 1
	}
	record, err := store.Load(id)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(record); err != nil {
		fmt.Fprintf(stderr, "write submission record: %v\n", err)
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
