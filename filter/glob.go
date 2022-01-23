package filter

import (
	"path/filepath"
	"strings"
)

func escapeBackSlash(pattern string) string {
	if filepath.Separator == '\\' {
		return strings.ReplaceAll(pattern, "\\", "\\\\")
	}
	return pattern
}
