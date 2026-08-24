package cosmicshortcut

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderDefaultPreservesExistingBindings(t *testing.T) {
	existing := "{\n    (modifiers: [Super], key: \"T\"): System(Terminal),\n}"
	got, err := RenderDefault(existing, "/usr/local/bin/cosmic-select")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "System(Terminal)") || !strings.Contains(got, "key: \"s\"") || !strings.Contains(got, "--activate") {
		t.Fatalf("shortcut was not preserved or added: %s", got)
	}
}

func TestRenderDefaultMigratesUppercaseShortcut(t *testing.T) {
	existing := "{\n    (modifiers: [Super, Shift], key: \"S\"): Spawn(\"/usr/local/bin/cosmic-select --activate\"),\n}\n"
	got, err := RenderDefault(existing, "/usr/local/bin/cosmic-select")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, `key: "S"`) || !strings.Contains(got, `key: "s"`) {
		t.Fatalf("uppercase shortcut was not migrated: %s", got)
	}
}

func TestInstallDefaultWritesExpectedPath(t *testing.T) {
	configDir := t.TempDir()
	path, err := InstallDefault(configDir, "/usr/local/bin/cosmic-select")
	if err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(configDir, configRelativePath)
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "Spawn(\"/usr/local/bin/cosmic-select --activate\")") {
		t.Fatalf("unexpected config: %s", content)
	}
}

func TestRenderDefaultRejectsMalformedConfig(t *testing.T) {
	if _, err := RenderDefault("not a RON map", "/usr/local/bin/cosmic-select"); err != ErrInvalidConfig {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}

func TestIsEphemeralExecutable(t *testing.T) {
	if !IsEphemeralExecutable(filepath.Join(os.TempDir(), "go-build123", "exe", "cosmic-select")) {
		t.Fatal("expected go run executable to be identified as ephemeral")
	}
	if IsEphemeralExecutable("/usr/local/bin/cosmic-select") {
		t.Fatal("expected installed executable to be persistent")
	}
}
