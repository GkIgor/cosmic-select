//go:build gtk4

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/GkIgor/cosmic-select/internal/portal"
	"github.com/GkIgor/cosmic-select/internal/ui"
)

func main() {
	client, err := portal.NewClient()
	if err != nil {
		ui.Run(os.Args, nil, err, "")
		return
	}
	defer client.Close()

	activationStatus := ""
	if hasActivation(os.Args) {
		activationStatus = activate(client)
	}
	ui.Run(os.Args, client, nil, activationStatus)
}

func hasActivation(args []string) bool {
	for _, arg := range args[1:] {
		if arg == "--activate" {
			return true
		}
	}
	return false
}

func activate(client *portal.Client) string {
	capabilities, err := client.CheckScreenshotCapabilities(context.Background())
	if err != nil {
		return fmt.Sprintf("Screenshot portal unavailable: %v", err)
	}
	selector := portal.NewScreenshotSelectorWithCapabilities(client, "", capabilities)
	image, err := selector.Select(context.Background())
	if err != nil {
		return fmt.Sprintf("Screen selection failed: %v", err)
	}
	return fmt.Sprintf("Screen region captured (%d bytes). Local OCR is the next step.", len(image.Data))
}
