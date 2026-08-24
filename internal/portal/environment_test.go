package portal

import "testing"

func TestEnvironmentRequiresCOSMICWayland(t *testing.T) {
	tests := []struct {
		name string
		env  Environment
		want bool
	}{
		{name: "cosmic wayland", env: Environment{Desktop: "COSMIC", Display: "wayland"}, want: true},
		{name: "cosmic in desktop list", env: Environment{Desktop: "foo:COSMIC", Display: "wayland"}, want: true},
		{name: "gnome wayland", env: Environment{Desktop: "GNOME", Display: "wayland"}, want: false},
		{name: "cosmic x11", env: Environment{Desktop: "COSMIC", Display: "x11"}, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.env.Validate() == nil
			if got != test.want {
				t.Fatalf("Validate() success = %v, want %v", got, test.want)
			}
		})
	}
}
