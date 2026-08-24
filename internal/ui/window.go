//go:build gtk4

package ui

import (
	"context"
	"fmt"
	"os"

	"github.com/GkIgor/cosmic-select/internal/portal"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const applicationID = "com.github.GkIgor.cosmicselect"

// Run starts the GTK application. Portal errors are shown in the window so
// the user gets a useful result even outside a supported COSMIC session.
func Run(args []string, client *portal.Client, startupErr error) {
	application := gtk.NewApplication(applicationID, gio.ApplicationFlagsNone)
	application.ConnectActivate(func() {
		showMainWindow(application, client, startupErr)
	})

	if code := application.Run(args); code > 0 {
		os.Exit(code)
	}
}

func showMainWindow(application *gtk.Application, client *portal.Client, startupErr error) {
	window := gtk.NewApplicationWindow(application)
	window.SetTitle("COSMIC Select")
	window.SetDefaultSize(420, 260)

	content := gtk.NewBox(gtk.OrientationVertical, 12)
	content.SetMarginTop(24)
	content.SetMarginBottom(24)
	content.SetMarginStart(24)
	content.SetMarginEnd(24)

	title := gtk.NewLabel("COSMIC Select")
	title.SetXalign(0)
	content.Append(title)

	status := gtk.NewLabel("Ready to connect to COSMIC portals.")
	status.SetXalign(0)
	status.SetWrap(true)
	content.Append(status)

	check := gtk.NewButtonWithLabel("Check portal support")
	check.ConnectClicked(func() {
		if startupErr != nil {
			status.SetText(fmt.Sprintf("Unable to connect to the session bus: %v", startupErr))
			return
		}
		if client == nil {
			status.SetText("Portal client is unavailable.")
			return
		}
		capabilities, err := client.CheckCapabilities(context.Background())
		if err != nil {
			status.SetText(fmt.Sprintf("COSMIC portal check failed: %v", err))
			return
		}
		status.SetText(fmt.Sprintf("COSMIC portals ready (Screenshot v%d, Global Shortcuts v%d).", capabilities.ScreenshotVersion, capabilities.GlobalShortcutsVersion))
	})
	content.Append(check)

	window.SetChild(content)
	window.Show()
}
