package config

import "testing"

func TestDefaults(t *testing.T) {
	settings := Defaults()
	if settings.GlobalShortcut == "" || settings.TranslationLanguage != "pt-BR" || settings.SearchEngine == "" {
		t.Fatalf("unexpected defaults: %+v", settings)
	}
}
