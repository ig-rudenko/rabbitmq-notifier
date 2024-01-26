package misc

import (
	"os"
	"strings"
)

func GetEnv(key, _default string) string {
	data := os.Getenv(key)
	if len(data) > 0 {
		return data
	}
	return _default
}

func JsonEscape(input string) string {
	escaped := input
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\n", "\\n")
	escaped = strings.ReplaceAll(escaped, "\r", "\\r")
	return escaped
}
