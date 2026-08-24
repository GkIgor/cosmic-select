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
	ScreenshotTargetOption bool
	GlobalShortcutsVersion uint32
}

func (c *Client) CheckCapabilities(ctx context.Context) (Capabilities, error) {
	if err := CurrentEnvironment().Validate(); err != nil {
		return Capabilities{}, err
	}
	screenshot, err := c.checkScreenshotCapabilities(ctx)
	if err != nil {
		return Capabilities{}, err
	}
	globalShortcutsVersion, err := c.globalShortcutsVersion(ctx)
	if err != nil {
		return Capabilities{}, err
	}
	screenshot.GlobalShortcutsVersion = globalShortcutsVersion
	return screenshot, nil
}

// CheckScreenshotCapabilities validates only the screenshot half of the
// integration. This remains useful when COSMIC lacks GlobalShortcuts.
func (c *Client) CheckScreenshotCapabilities(ctx context.Context) (Capabilities, error) {
	if err := CurrentEnvironment().Validate(); err != nil {
		return Capabilities{}, err
	}
	return c.checkScreenshotCapabilities(ctx)
}

func (c *Client) checkScreenshotCapabilities(ctx context.Context) (Capabilities, error) {
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
	screenshotTargets, targetsErr := uint32Property(object, screenshotTargetsProperty)
	screenshotTargetOption, err := resolveScreenshotTargetOption(screenshotVersion, screenshotTargets, targetsErr)
	if err != nil {
		return Capabilities{}, err
	}
	if !screenshotTargetOption && targetsErr == nil {
		return Capabilities{}, fmt.Errorf("%w: area selection is not supported", ErrPortalUnavailable)
	}

	return Capabilities{
		ScreenshotVersion:      screenshotVersion,
		ScreenshotTargets:      screenshotTargets,
		ScreenshotTargetOption: screenshotTargetOption,
	}, nil
}

func (c *Client) globalShortcutsVersion(ctx context.Context) (uint32, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	version, err := uint32Property(c.object(), globalShortcutsVersionProperty)
	if err != nil {
		return 0, fmt.Errorf("%w: global shortcuts portal: %v", ErrPortalUnavailable, err)
	}
	return version, nil
}

// resolveScreenshotTargetOption supports Screenshot portal v2, which exposes
// interactive selection but does not expose the v3 AvailableTargets property.
func resolveScreenshotTargetOption(version, targets uint32, targetsErr error) (bool, error) {
	if targetsErr != nil {
		if version >= 2 {
			return false, nil
		}
		return false, fmt.Errorf("%w: interactive screenshot is unavailable", ErrPortalUnavailable)
	}
	if targets&areaTarget == 0 {
		return false, fmt.Errorf("%w: area selection is not supported", ErrPortalUnavailable)
	}
	return true, nil
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
