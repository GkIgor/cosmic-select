package ocr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/GkIgor/cosmic-select/internal/ports"
)

var ErrInvalidImage = errors.New("OCR received an empty image")

type commandRunner interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

type execCommandRunner struct{}

func (execCommandRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).Output()
}

// TesseractEngine runs Tesseract locally against a transient PNG file.
type TesseractEngine struct {
	binary    string
	languages []string
	runner    commandRunner
}

func NewTesseractEngine(binary string, languages []string) *TesseractEngine {
	if binary == "" {
		binary = "tesseract"
	}
	if len(languages) == 0 {
		languages = []string{"eng", "por"}
	}
	return &TesseractEngine{
		binary:    binary,
		languages: append([]string(nil), languages...),
		runner:    execCommandRunner{},
	}
}

func newTesseractEngineWithRunner(binary string, languages []string, runner commandRunner) *TesseractEngine {
	engine := NewTesseractEngine(binary, languages)
	engine.runner = runner
	return engine
}

func (e *TesseractEngine) Extract(ctx context.Context, image ports.Image) (string, error) {
	if len(image.Data) == 0 {
		return "", ErrInvalidImage
	}

	temporary, err := os.CreateTemp("", "cosmic-select-ocr-*.png")
	if err != nil {
		return "", fmt.Errorf("create transient OCR image: %w", err)
	}
	path := temporary.Name()
	defer os.Remove(path)
	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		return "", fmt.Errorf("protect transient OCR image: %w", err)
	}
	if _, err := temporary.Write(image.Data); err != nil {
		_ = temporary.Close()
		return "", fmt.Errorf("write transient OCR image: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return "", fmt.Errorf("close transient OCR image: %w", err)
	}

	args := []string{path, "stdout", "-l", strings.Join(e.languages, "+"), "--psm", "6"}
	output, err := e.runner.Run(ctx, e.binary, args...)
	if err != nil {
		return "", fmt.Errorf("run local OCR: %w", err)
	}
	return CleanText(string(output)), nil
}
