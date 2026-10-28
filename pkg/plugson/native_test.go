package plugson

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNativePlugson(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "plugson_native_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	np := NewNativePlugson()
	np.Config.ThemeName = "custom-adwaita"
	np.Config.Timeout = 15

	if err := np.SaveConfig(tmpDir); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	jsonPath := filepath.Join(tmpDir, "ventoy.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("failed to read ventoy.json: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("ventoy.json is empty")
	}
}
