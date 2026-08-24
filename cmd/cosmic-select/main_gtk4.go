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
		ui.Run(os.Args, nil, err)
		return
	}
	defer client.Close()
	ui.Run(os.Args, client, nil)
}
