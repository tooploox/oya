package shell

import (
	"bufio"
	"fmt"
	"strings"
)

func firstLine(s string) string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	if scanner.Scan() && scanner.Err() == nil {
		return scanner.Text()
	} else {
		return ""
	}
}

func shorten(s string, max int) string {
	if len(s) > max {
		return fmt.Sprintf("%s...", s[:max-3])
	} else {
		return s
	}
}
