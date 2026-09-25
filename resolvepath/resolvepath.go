// Package resolvepath provides a shared way for jcli tools to determine
// their primary input path, so they can be chained together the Unix way.
package resolvepath

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Resolve determines the manuscript path to operate on. It prefers an
// explicit argv[1], and falls back to reading a single line from stdin so
// that jcli tools can be chained together with a plain Unix pipe, e.g.
// `jcli-upload manuscript.md | jcli-validate`.
func Resolve(args []string, stdin io.Reader) (string, error) {
	if len(args) >= 2 && strings.TrimSpace(args[1]) != "" {
		return args[1], nil
	}

	scanner := bufio.NewScanner(stdin)
	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			return line, nil
		}
	}

	return "", fmt.Errorf("no manuscript path provided as an argument or via stdin")
}
