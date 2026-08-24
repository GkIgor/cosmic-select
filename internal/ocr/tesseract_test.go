package ocr

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/GkIgor/cosmic-select/internal/ports"
)

type runnerStub struct {
	name string
	args []string
	text []byte
}

func (r *runnerStub) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	r.name = name
	r.args = args
	if len(args) == 0 {
		return nil, nil
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return nil, err
	}
	if string(data) != "png bytes" {
		return nil, os.ErrInvalid
	}
	return r.text, nil
}

func TestTesseractExtractsAndCleansText(t *testing.T) {
	runner := &runnerStub{text: []byte("  Linha um  \n\n Linha dois \n")}
	engine := newTesseractEngineWithRunner("tesseract", []string{"eng", "por"}, runner)

	got, err := engine.Extract(context.Background(), ports.Image{Data: []byte("png bytes")})
	if err != nil {
		t.Fatal(err)
	}
	if got != "Linha um\nLinha dois" {
		t.Fatalf("text = %q", got)
	}
	if runner.name != "tesseract" || strings.Join(runner.args[2:], " ") != "-l eng+por --psm 6" {
		t.Fatalf("unexpected Tesseract invocation: %s %v", runner.name, runner.args)
	}
	if _, err := os.Stat(runner.args[0]); !os.IsNotExist(err) {
		t.Fatalf("temporary OCR image still exists: %v", err)
	}
}

func TestTesseractRejectsEmptyImage(t *testing.T) {
	engine := NewTesseractEngine("tesseract", nil)
	if _, err := engine.Extract(context.Background(), ports.Image{}); err != ErrInvalidImage {
		t.Fatalf("expected ErrInvalidImage, got %v", err)
	}
}
