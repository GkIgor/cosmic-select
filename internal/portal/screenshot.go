package portal

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/GkIgor/cosmic-select/internal/ports"
	"github.com/godbus/dbus/v5"
)

const screenshotInterface = "org.freedesktop.portal.Screenshot"

// ScreenshotSelector delegates interactive area selection to the Screenshot
// portal. The returned bytes are transient and belong to the caller.
type ScreenshotSelector struct {
	client       *Client
	parentWindow string
}

func NewScreenshotSelector(client *Client, parentWindow string) *ScreenshotSelector {
	return &ScreenshotSelector{client: client, parentWindow: parentWindow}
}

func (s *ScreenshotSelector) Select(ctx context.Context) (ports.Image, error) {
	s.client.mu.Lock()
	defer s.client.mu.Unlock()

	call := s.client.object().CallWithContext(ctx, screenshotInterface+".Screenshot", 0,
		s.parentWindow,
		map[string]dbus.Variant{
			"handle_token": dbus.MakeVariant("cosmic_select_screenshot"),
			"interactive":  dbus.MakeVariant(true),
			"target":       dbus.MakeVariant(uint32(4)),
		},
	)
	if call.Err != nil {
		return ports.Image{}, fmt.Errorf("request screenshot: %w", call.Err)
	}

	var request dbus.ObjectPath
	if err := call.Store(&request); err != nil {
		return ports.Image{}, fmt.Errorf("read screenshot request handle: %w", err)
	}
	results, err := s.client.waitResponse(ctx, request)
	if err != nil {
		return ports.Image{}, err
	}
	uri, err := variantString(results, "uri")
	if err != nil {
		return ports.Image{}, err
	}
	return readScreenshotURI(uri)
}

func readScreenshotURI(rawURI string) (ports.Image, error) {
	parsed, err := url.Parse(rawURI)
	if err != nil || parsed.Scheme != "file" {
		return ports.Image{}, ErrUnsupportedScreenshot
	}
	path, err := url.PathUnescape(parsed.Path)
	if err != nil || filepath.Clean(path) != path {
		return ports.Image{}, ErrUnsupportedScreenshot
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ports.Image{}, fmt.Errorf("read portal screenshot: %w", err)
	}
	return ports.Image{Data: data, MimeType: "image/png"}, nil
}
