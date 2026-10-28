package theme

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestThemeMaker(t *testing.T) {
	cfg := ThemeConfig{
		Title:        "Custom Multiverse Boot",
		PrimaryColor: "#2563eb",
	}

	tm := NewThemeMaker(cfg)
	txt := tm.GenerateThemeTxt()

	if !strings.Contains(txt, "Custom Multiverse Boot") {
		t.Errorf("expected theme title in theme.txt, got:\n%s", txt)
	}

	tmpDir, err := os.MkdirTemp("", "theme_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := tm.SaveTheme(tmpDir); err != nil {
		t.Fatalf("SaveTheme failed: %v", err)
	}

	themePath := filepath.Join(tmpDir, "theme.txt")
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		t.Errorf("theme.txt was not created")
	}
}
