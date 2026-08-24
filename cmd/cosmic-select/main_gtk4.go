//go:build gtk4

package main

import (
	"os"

	"github.com/GkIgor/cosmic-select/internal/portal"
	"github.com/GkIgor/cosmic-select/internal/ui"
)

func main() {
	client, err := portal.NewClient()
	if err != nil {
		ui.Run(os.Args, nil, err, hasActivation(os.Args))
		return
	}
	defer client.Close()
	ui.Run(os.Args, client, nil, hasActivation(os.Args))
}

func hasActivation(args []string) bool {
	for _, arg := range args[1:] {
		if arg == "--activate" {
			return true
		}
	}
	return false
}
