package portal

import (
	"os"
	"strings"
)

// Environment contains the desktop assumptions required by COSMIC Select.
type Environment struct {
	Desktop string
	Display string
}

func CurrentEnvironment() Environment {
	return Environment{
		Desktop: os.Getenv("XDG_CURRENT_DESKTOP"),
		Display: os.Getenv("XDG_SESSION_TYPE"),
	}
}

func (e Environment) Validate() error {
	if !containsDesktop(e.Desktop, "cosmic") || !strings.EqualFold(e.Display, "wayland") {
		return ErrUnsupportedEnvironment
	}
	return nil
}

func containsDesktop(value, expected string) bool {
	for _, desktop := range strings.FieldsFunc(value, func(r rune) bool { return r == ':' || r == ';' }) {
		if strings.EqualFold(strings.TrimSpace(desktop), expected) {
			return true
		}
	}
	return false
}
