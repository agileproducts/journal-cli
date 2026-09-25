// Package resolvepath provides a shared way for jcli tools to determine
// their primary input (a path or submission ID) for Unix-style chaining.
package resolvepath

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Resolve determines the input path or submission ID to operate on. It prefers an
// explicit argv[1], and falls back to reading a single line from stdin so
// that jcli tools can be chained together with a plain Unix pipe, e.g.
// `jcli-upload manuscript.md | jcli-validate`.
func Resolve(args []string, stdin io.Reader) (string, error) {
	if len(args) >= 2 {
		if strings.TrimSpace(args[1]) == "" {
			return "", fmt.Errorf("input argument must not be empty")
		}
		return args[1], nil
	}

	scanner := bufio.NewScanner(stdin)
	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}
	return "", fmt.Errorf("no input provided as an argument or via stdin")
}
