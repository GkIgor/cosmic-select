package portal

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	screenshotVersionProperty      = screenshotInterface + ".version"
	screenshotTargetsProperty      = screenshotInterface + ".AvailableTargets"
	globalShortcutsVersionProperty = globalShortcutsInterface + ".version"
	areaTarget                     = uint32(4)
)

// Capabilities describes the portal features required by COSMIC Select.
type Capabilities struct {
	ScreenshotVersion      uint32
	ScreenshotTargets      uint32
	GlobalShortcutsVersion uint32
}

func (c *Client) CheckCapabilities(ctx context.Context) (Capabilities, error) {
	if err := CurrentEnvironment().Validate(); err != nil {
		return Capabilities{}, err
	}
	select {
	case <-ctx.Done():
		return Capabilities{}, ctx.Err()
	default:
	}

	object := c.object()
	screenshotVersion, err := uint32Property(object, screenshotVersionProperty)
	if err != nil {
		return Capabilities{}, fmt.Errorf("%w: screenshot portal: %v", ErrPortalUnavailable, err)
	}
	screenshotTargets, err := uint32Property(object, screenshotTargetsProperty)
	if err != nil {
		return Capabilities{}, fmt.Errorf("%w: screenshot targets: %v", ErrPortalUnavailable, err)
	}
	globalShortcutsVersion, err := uint32Property(object, globalShortcutsVersionProperty)
	if err != nil {
		return Capabilities{}, fmt.Errorf("%w: global shortcuts portal: %v", ErrPortalUnavailable, err)
	}
	if screenshotTargets&areaTarget == 0 {
		return Capabilities{}, fmt.Errorf("%w: area selection is not supported", ErrPortalUnavailable)
	}

	return Capabilities{
		ScreenshotVersion:      screenshotVersion,
		ScreenshotTargets:      screenshotTargets,
		GlobalShortcutsVersion: globalShortcutsVersion,
	}, nil
}

func uint32Property(object dbus.BusObject, property string) (uint32, error) {
	value, err := object.GetProperty(property)
	if err != nil {
		return 0, err
	}
	result, ok := value.Value().(uint32)
	if !ok {
		return 0, fmt.Errorf("property %q has unexpected type", property)
	}
	return result, nil
}
