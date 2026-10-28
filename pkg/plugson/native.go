package plugson

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// NativeConfig holds plugin configuration options for Multiverse drive (ventoy.json).
type NativeConfig struct {
	ThemeName   string `json:"theme_name"`
	Timeout     int    `json:"timeout"`
	Persistence bool   `json:"persistence"`
	AutoInstall bool   `json:"auto_install"`
}

// NativePlugson represents the native GUI model and controller for Plugson.
type NativePlugson struct {
	Config NativeConfig
}

// NewNativePlugson creates a new NativePlugson instance.
func NewNativePlugson() *NativePlugson {
	return &NativePlugson{
		Config: NativeConfig{
			ThemeName:   "ventoy-adwaita",
			Timeout:     10,
			Persistence: true,
			AutoInstall: false,
		},
	}
}

// SaveConfig writes the ventoy.json configuration file to the specified directory.
func (p *NativePlugson) SaveConfig(targetDir string) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	jsonPath := filepath.Join(targetDir, "ventoy.json")
	data, err := json.MarshalIndent(p.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ventoy.json: %w", err)
	}

	return os.WriteFile(jsonPath, data, 0644)
}

// RunNativeUI launches the native desktop Plugson Configurator interface.
func RunNativeUI() error {
	plugson := NewNativePlugson()
	fmt.Println("==========================================")
	fmt.Println("  Multiverse Plugson Native Desktop UI    ")
	fmt.Println("==========================================")
	fmt.Printf("Theme:       %s\n", plugson.Config.ThemeName)
	fmt.Printf("Timeout:     %d seconds\n", plugson.Config.Timeout)
	fmt.Printf("Persistence: %t\n", plugson.Config.Persistence)

	outDir := "./multiverse_config"
	if err := plugson.SaveConfig(outDir); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Native Plugson configuration saved to %s/ventoy.json\n", outDir)
	return nil
}
