package ocr

import "testing"

func TestCleanText(t *testing.T) {
	input := "  First line\r\n\r\n  Second line  \n"
	if got, want := CleanText(input), "First line\nSecond line"; got != want {
		t.Fatalf("CleanText() = %q, want %q", got, want)
	}
}

func TestCleanTextEmpty(t *testing.T) {
	if got := CleanText(" \n\t "); got != "" {
		t.Fatalf("CleanText() = %q, want empty text", got)
	}
}
