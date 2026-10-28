package theme

import (
	"fmt"
	"os"
	"path/filepath"
)

// ThemeConfig holds parameters for constructing a customized GRUB Adwaita theme.
type ThemeConfig struct {
	Title           string `json:"title"`
	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	FontName        string `json:"font_name"`
	Resolution      string `json:"resolution"`
	ShowProgressBar bool   `json:"show_progress_bar"`
}

// ThemeMaker handles generating GRUB theme configurations and asset bundles.
type ThemeMaker struct {
	config ThemeConfig
}

// NewThemeMaker creates a new ThemeMaker instance.
func NewThemeMaker(cfg ThemeConfig) *ThemeMaker {
	if cfg.Title == "" {
		cfg.Title = "Multiverse Boot Menu"
	}
	if cfg.PrimaryColor == "" {
		cfg.PrimaryColor = "#2563eb"
	}
	if cfg.SecondaryColor == "" {
		cfg.SecondaryColor = "#cbd5e1"
	}
	if cfg.Resolution == "" {
		cfg.Resolution = "1024x768"
	}
	return &ThemeMaker{config: cfg}
}

// GenerateThemeTxt generates the GRUB theme.txt content formatted with Adwaita styling.
func (tm *ThemeMaker) GenerateThemeTxt() string {
	return fmt.Sprintf(`# Multiverse Adwaita Theme Configuration
desktop-image: "background.png"
title-text: "%s"
title-color: "%s"
title-font: "Unifont Regular 16"
desktop-color: "%s"

+ boot_menu {
  left = 15%%
  width = 70%%
  top = 20%%
  height = 60%%
  item_color = "#ffffff"
  selected_item_color = "%s"
  icon_width = 32
  icon_height = 32
  item_font = "Unifont Regular 14"
}

+ progress_bar {
  id = "__timeout__"
  left = 15%%
  width = 70%%
  top = 85%%
  height = 16
  bg_color = "%s"
  fg_color = "%s"
}
`, tm.config.Title, tm.config.PrimaryColor, tm.config.SecondaryColor, tm.config.PrimaryColor, tm.config.SecondaryColor, tm.config.PrimaryColor)
}

// SaveTheme writes theme.txt and directory assets to target directory.
func (tm *ThemeMaker) SaveTheme(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create theme output directory: %w", err)
	}

	themeTxtPath := filepath.Join(outputDir, "theme.txt")
	content := tm.GenerateThemeTxt()

	return os.WriteFile(themeTxtPath, []byte(content), 0644)
}
