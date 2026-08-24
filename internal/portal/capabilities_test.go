package portal

import (
	"errors"
	"testing"
)

func TestResolveScreenshotTargetOption(t *testing.T) {
	propertyErr := errors.New("Unknown property AvailableTargets")

	tests := []struct {
		name      string
		version   uint32
		targets   uint32
		err       error
		useTarget bool
		wantErr   bool
	}{
		{name: "v3 area target", version: 3, targets: areaTarget, useTarget: true},
		{name: "v3 without area target", version: 3, targets: 1, wantErr: true},
		{name: "v2 interactive fallback", version: 2, err: propertyErr},
		{name: "v1 without target property", version: 1, err: propertyErr, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useTarget, err := resolveScreenshotTargetOption(test.version, test.targets, test.err)
			if (err != nil) != test.wantErr {
				t.Fatalf("error = %v, want error: %v", err, test.wantErr)
			}
			if err == nil && useTarget != test.useTarget {
				t.Fatalf("useTarget = %v, want %v", useTarget, test.useTarget)
			}
		})
	}
}
