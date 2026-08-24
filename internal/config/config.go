package config

const defaultShortcut = "<Super><Shift>s"

// Settings contains only preferences required by the core flow.
type Settings struct {
	GlobalShortcut      string
	TranslationLanguage string
	SearchEngine        string
}

func Defaults() Settings {
	return Settings{
		GlobalShortcut:      defaultShortcut,
		TranslationLanguage: "pt-BR",
		SearchEngine:        "https://www.google.com/search?q=%s",
	}
}
