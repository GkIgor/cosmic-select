package ocr

import (
	"strings"
	"unicode"
)

// CleanText normalizes OCR output before it reaches an action provider.
// It preserves line breaks while removing accidental surrounding whitespace.
func CleanText(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	cleaned := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimFunc(line, unicode.IsSpace)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}
