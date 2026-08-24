package ports

import "context"

// Image is a transient screenshot selected through the desktop portal.
// Implementations must not persist it beyond the active operation.
type Image struct {
	Data     []byte
	MimeType string
}

// Selector delegates screen selection and capture to the COSMIC portal.
type Selector interface {
	Select(ctx context.Context) (Image, error)
}

// Shortcut registers the application activation shortcut through the XDG
// Global Shortcuts portal.
type Shortcut interface {
	Register(ctx context.Context, trigger func()) error
	Close() error
}

// OCREngine runs locally and returns text from a transient selected image.
type OCREngine interface {
	Extract(ctx context.Context, image Image) (string, error)
}
