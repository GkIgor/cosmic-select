//go:build gtk4

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/GkIgor/cosmic-select/internal/ocr"
	"github.com/GkIgor/cosmic-select/internal/portal"
	"github.com/GkIgor/cosmic-select/internal/ports"
	"github.com/GkIgor/cosmic-select/internal/ui"
)

func main() {
	client, err := portal.NewClient()
	if err != nil {
		ui.Run(os.Args, nil, err, "", "")
		return
	}
	defer client.Close()

	activationStatus := ""
	recognizedText := ""
	if hasActivation(os.Args) {
		activationStatus, recognizedText = activate(client)
	}
	ui.Run(os.Args, client, nil, activationStatus, recognizedText)
}

func hasActivation(args []string) bool {
	for _, arg := range args[1:] {
		if arg == "--activate" {
			return true
		}
	}
	return false
}

func activate(client *portal.Client) (string, string) {
	capabilities, err := client.CheckScreenshotCapabilities(context.Background())
	if err != nil {
		return fmt.Sprintf("Screenshot portal unavailable: %v", err), ""
	}
	selector := portal.NewScreenshotSelectorWithCapabilities(client, "", capabilities)
	image, err := selector.Select(context.Background())
	if err != nil {
		return fmt.Sprintf("Screen selection failed: %v", err), ""
	}
	engine := ocr.NewTesseractEngine("tesseract", []string{"eng", "por"})
	text, err := engine.Extract(context.Background(), ports.Image{Data: image.Data, MimeType: image.MimeType})
	if err != nil {
		return fmt.Sprintf("Local OCR failed: %v", err), ""
	}
	if text == "" {
		return "No text detected.", ""
	}
	return "Local OCR completed.", text
}
