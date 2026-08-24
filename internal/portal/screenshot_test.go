package portal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadScreenshotURI(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "capture.png")
	if err := os.WriteFile(path, []byte("png bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	image, err := readScreenshotURI("file://" + path)
	if err != nil {
		t.Fatalf("readScreenshotURI() error = %v", err)
	}
	if string(image.Data) != "png bytes" || image.MimeType != "image/png" {
		t.Fatalf("unexpected image: %+v", image)
	}
}

func TestReadScreenshotURIRejectsNonFileURI(t *testing.T) {
	if _, err := readScreenshotURI("https://example.test/capture.png"); err != ErrUnsupportedScreenshot {
		t.Fatalf("expected ErrUnsupportedScreenshot, got %v", err)
	}
}
