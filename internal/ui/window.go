//go:build gtk4

package ui

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/GkIgor/cosmic-select/internal/cosmicshortcut"
	"github.com/GkIgor/cosmic-select/internal/portal"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const applicationID = "com.github.GkIgor.cosmicselect"

// Run starts the GTK application. Portal errors are shown in the window so
// the user gets a useful result even outside a supported COSMIC session.
func Run(args []string, client *portal.Client, startupErr error, activationStatus string) {
	application := gtk.NewApplication(applicationID, gio.ApplicationFlagsNone)
	application.ConnectActivate(func() {
		showMainWindow(application, client, startupErr, activationStatus)
	})

	if code := application.Run(args); code > 0 {
		os.Exit(code)
	}
}

func showMainWindow(application *gtk.Application, client *portal.Client, startupErr error, activationStatus string) {
	window := gtk.NewApplicationWindow(application)
	window.SetTitle("COSMIC Select")
	window.SetDefaultSize(420, 260)

	content := gtk.NewBox(gtk.OrientationVertical, 12)
	content.SetMarginTop(24)
	content.SetMarginBottom(24)
	content.SetMarginStart(24)
	content.SetMarginEnd(24)

	title := gtk.NewLabel("COSMIC Select")
	title.SetXAlign(0)
	content.Append(title)

	initialStatus := "Ready to connect to COSMIC portals."
	if activationStatus != "" {
		initialStatus = activationStatus
	}
	status := gtk.NewLabel(initialStatus)
	status.SetXAlign(0)
	status.SetWrap(true)
	content.Append(status)

	check := gtk.NewButtonWithLabel("Check portal support")
	install := gtk.NewButtonWithLabel("Install COSMIC shortcut (Super + Shift + S)")
	install.SetVisible(false)
	check.ConnectClicked(func() {
		if startupErr != nil {
			status.SetText(fmt.Sprintf("Unable to connect to the session bus: %v", startupErr))
			install.SetVisible(true)
			return
		}
		if client == nil {
			status.SetText("Portal client is unavailable.")
			return
		}
		capabilities, err := client.CheckCapabilities(context.Background())
		if err != nil {
			status.SetText(fmt.Sprintf("COSMIC portal check failed: %v", err))
			install.SetVisible(errors.Is(err, portal.ErrPortalUnavailable))
			return
		}
		install.SetVisible(false)
		status.SetText(fmt.Sprintf("COSMIC portals ready (Screenshot v%d, Global Shortcuts v%d).", capabilities.ScreenshotVersion, capabilities.GlobalShortcutsVersion))
	})
	content.Append(check)

	install.ConnectClicked(func() {
		executable, err := os.Executable()
		if err != nil {
			status.SetText(fmt.Sprintf("Unable to find COSMIC Select executable: %v", err))
			return
		}
		if cosmicshortcut.IsEphemeralExecutable(executable) {
			status.SetText("Build a persistent binary before installing the shortcut: go build -tags gtk4 -o $HOME/.local/bin/cosmic-select ./cmd/cosmic-select")
			return
		}
		configDir, err := os.UserConfigDir()
		if err != nil {
			status.SetText(fmt.Sprintf("Unable to find COSMIC configuration directory: %v", err))
			return
		}
		path, err := cosmicshortcut.InstallDefault(configDir, executable)
		if err != nil {
			status.SetText(fmt.Sprintf("Unable to install native COSMIC shortcut: %v", err))
			return
		}
		status.SetText(fmt.Sprintf("COSMIC shortcut installed at %s. Restart COSMIC if it does not activate immediately.", path))
		install.SetVisible(false)
	})
	content.Append(install)

	window.SetChild(content)
	window.Show()
}
