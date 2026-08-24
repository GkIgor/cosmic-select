package portal

import "errors"

var (
	ErrUnsupportedEnvironment = errors.New("COSMIC Desktop on Wayland is required")
	ErrPortalUnavailable      = errors.New("required desktop portal is unavailable")
	ErrPortalCancelled        = errors.New("portal request was cancelled")
	ErrPortalResponse         = errors.New("portal request failed")
	ErrUnsupportedScreenshot  = errors.New("portal returned an unsupported screenshot URI")
)
