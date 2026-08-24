package cosmicshortcut

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var ErrInvalidConfig = errors.New("invalid COSMIC shortcut configuration")

const (
	defaultKey         = "s"
	configRelativePath = "com.system76.CosmicSettings.Shortcuts/v1/custom"
)

// DefaultCommand returns the executable command used by COSMIC's native
// shortcut action. The executable must be a persistent installed binary.
func DefaultCommand(executable string) string {
	return fmt.Sprintf("Spawn(%s)", strconv.Quote(executable+" --activate"))
}

func ConfigPath(configDir string) string {
	return filepath.Join(configDir, configRelativePath)
}

func IsEphemeralExecutable(executable string) bool {
	relative, err := filepath.Rel(os.TempDir(), filepath.Clean(executable))
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return false
	}
	firstComponent := strings.Split(relative, string(os.PathSeparator))[0]
	return strings.HasPrefix(firstComponent, "go-build")
}

// RenderDefault appends the COSMIC-native shortcut while preserving existing
// custom bindings. COSMIC stores this file as a small RON map.
func RenderDefault(existing, executable string) (string, error) {
	trimmed := strings.TrimSpace(existing)
	if trimmed == "" {
		return renderMap("", executable), nil
	}
	if !strings.HasPrefix(trimmed, "{") || !strings.HasSuffix(trimmed, "}") {
		return "", ErrInvalidConfig
	}
	body := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	currentEntry := defaultEntry(executable)
	if strings.Contains(trimmed, currentEntry) {
		return existing, nil
	}
	legacyEntry := fmt.Sprintf("    (modifiers: [Super, Shift], key: %q): %s,", "S", DefaultCommand(executable))
	if strings.Contains(existing, legacyEntry) {
		return strings.Replace(existing, legacyEntry, currentEntry, 1), nil
	}
	return renderMap(body, executable), nil
}

// InstallDefault writes the native shortcut atomically into configDir.
func InstallDefault(configDir, executable string) (string, error) {
	if executable == "" {
		return "", errors.New("COSMIC Select executable path is required")
	}

	path := ConfigPath(configDir)
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", fmt.Errorf("create COSMIC shortcut directory: %w", err)
	}

	existing, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("read COSMIC shortcut configuration: %w", err)
	}
	content, err := RenderDefault(string(existing), executable)
	if err != nil {
		return "", err
	}
	if content == string(existing) {
		return path, nil
	}

	temporary, err := os.CreateTemp(directory, ".custom.tmp-*")
	if err != nil {
		return "", fmt.Errorf("create temporary shortcut configuration: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		return "", fmt.Errorf("set shortcut configuration permissions: %w", err)
	}
	if _, err := temporary.WriteString(content); err != nil {
		_ = temporary.Close()
		return "", fmt.Errorf("write shortcut configuration: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return "", fmt.Errorf("close shortcut configuration: %w", err)
	}
	if err := os.Rename(temporaryName, path); err != nil {
		return "", fmt.Errorf("install COSMIC shortcut configuration: %w", err)
	}
	return path, nil
}

func renderMap(existing, executable string) string {
	entry := defaultEntry(executable)
	if existing == "" {
		return "{\n" + entry + "}\n"
	}
	if strings.HasSuffix(existing, ",") {
		return "{\n" + existing + "\n" + entry + "}\n"
	}
	return "{\n" + existing + ",\n" + entry + "}\n"
}

func defaultEntry(executable string) string {
	return fmt.Sprintf("    (modifiers: [Super, Shift], key: %q): %s,\n", defaultKey, DefaultCommand(executable))
}
