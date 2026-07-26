package execute

import (
	"fmt"
	"strings"
)

func buildSourceDiagnostic(path, source string, line, column int, message string) string {
	if line < 1 {
		line = 1
	}
	if column < 1 {
		column = 1
	}

	lines := strings.Split(source, "\n")
	lineText := ""
	if line-1 >= 0 && line-1 < len(lines) {
		lineText = lines[line-1]
	}

	caretPadding := strings.Repeat(" ", column-1)
	return fmt.Sprintf("%s:%d:%d: error: %s\n  | %s\n  | %s^", path, line, column, message, lineText, caretPadding)
}
