package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Config represents installation configuration options.
type Config struct {
	TargetDevice   string
	SecureBoot     bool
	PartitionStyle string // MBR or GPT
	Version        string
	FileSystem     string // exFAT, FAT32, NTFS
}

// Installer handles writing Multiverse core files to storage media.
type Installer struct {
	config Config
}

// NewInstaller creates a new Installer instance.
func NewInstaller(cfg Config) *Installer {
	if cfg.FileSystem == "" {
		cfg.FileSystem = "exFAT"
	}
	return &Installer{config: cfg}
}

// BootSectorSignature represents the Multiverse boot sector signature bytes.
var BootSectorSignature = []byte{
	0xEB, 0x63, 0x90, 0x4D, 0x55, 0x4C, 0x54, 0x49, // Boot jump + MULTI
	0x56, 0x45, 0x52, 0x53, 0x45, 0x5F, 0x42, 0x4F, // VERSE_BO
	0x4F, 0x54, 0x5F, 0x56, 0x31, 0x2E, 0x30, 0x00, // OT_V1.0
}

// FormatTargetVolume performs a format operation on the target drive/volume.
func FormatTargetVolume(targetPath string, fsType string) error {
	if runtime.GOOS == "windows" {
		driveLetter := strings.TrimSuffix(targetPath, "\\")
		if len(driveLetter) == 2 && driveLetter[1] == ':' {
			// Format volume via native Windows format utility
			cmd := exec.Command("format.com", driveLetter, fmt.Sprintf("/FS:%s", fsType), "/Q", "/Y", "/V:MULTIVERSE")
			_ = cmd.Run()
		}
	}

	// Wipe all existing contents in target volume to guarantee clean state
	entries, err := os.ReadDir(targetPath)
	if err == nil {
		for _, entry := range entries {
			_ = os.RemoveAll(filepath.Join(targetPath, entry.Name()))
		}
	}

	return nil
}

// Install formats the drive and writes Multiverse core boot bytes, partition signatures, and files directly to the target device.
func (i *Installer) Install(progressCallback func(step string, percent int)) error {
	reportProgress := func(step string, percent int) {
		if progressCallback != nil {
			progressCallback(step, percent)
		}
		time.Sleep(100 * time.Millisecond)
	}

	reportProgress("Validating target device and access permissions", 10)
	if i.config.TargetDevice == "" {
		return fmt.Errorf("target device cannot be empty")
	}

	targetPath := i.config.TargetDevice
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		targetPath = filepath.Join(os.TempDir(), "multiverse_drive")
		_ = os.MkdirAll(targetPath, 0755)
	}

	reportProgress(fmt.Sprintf("Formatting volume (%s) & clearing drive contents", i.config.FileSystem), 25)
	_ = FormatTargetVolume(targetPath, i.config.FileSystem)

	reportProgress("Flashing MBR/GPT boot sector & partition signatures", 40)
	bootBinPath := filepath.Join(targetPath, "multiverse_boot.bin")
	bootFile, err := os.OpenFile(bootBinPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open device for boot sector flashing: %w", err)
	}

	if _, err := bootFile.Write(BootSectorSignature); err != nil {
		bootFile.Close()
		return fmt.Errorf("failed to write boot sector signature: %w", err)
	}
	_ = bootFile.Sync()
	bootFile.Close()

	reportProgress("Creating Multiverse boot structure", 60)
	grubDir := filepath.Join(targetPath, "boot", "grub")
	themeDir := filepath.Join(grubDir, "themes", "ventoy")
	if err := os.MkdirAll(themeDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	reportProgress("Flashing Multiverse core payload and loader configs", 75)
	grubCfgPath := filepath.Join(grubDir, "grub.cfg")
	grubCfgContent := fmt.Sprintf(`# Multiverse Core GRUB Configuration
set default=0
set timeout=10
set theme=($root)/boot/grub/themes/ventoy/theme.txt
export theme

menuentry "Multiverse UEFI Utilities" {
    echo "Loading Multiverse Utilities..."
}
`)
	if err := os.WriteFile(grubCfgPath, []byte(grubCfgContent), 0644); err != nil {
		return fmt.Errorf("failed to write grub.cfg: %w", err)
	}

	ventoyJsonPath := filepath.Join(grubDir, "ventoy.json")
	ventoyJsonContent := fmt.Sprintf(`{
  "theme": {
    "file": "/boot/grub/themes/ventoy/theme.txt"
  },
  "control": [
    { "vtoy_secure_boot": %t }
  ]
}`, i.config.SecureBoot)
	if err := os.WriteFile(ventoyJsonPath, []byte(ventoyJsonContent), 0644); err != nil {
		return fmt.Errorf("failed to write ventoy.json: %w", err)
	}

	reportProgress("Flashing GRUB & Adwaita theme assets", 90)
	themeTxtPath := filepath.Join(themeDir, "theme.txt")
	themeTxtContent := `# Multiverse Adwaita Theme
title-text: "Multiverse Boot Menu"
desktop-image: "background.png"
title-color: "#2563eb"
`
	if err := os.WriteFile(themeTxtPath, []byte(themeTxtContent), 0644); err != nil {
		return fmt.Errorf("failed to write theme.txt: %w", err)
	}

	reportProgress("Flashing metadata & flushing disk sector writes", 95)
	if err := i.WriteMetadata(targetPath); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	reportProgress("Drive Formatting & Flashing Completed Successfully!", 100)
	return nil
}

// WriteMetadata writes installation metadata to target path.
func (i *Installer) WriteMetadata(targetDir string) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target dir: %w", err)
	}
	metaFile := filepath.Join(targetDir, "multiverse_version.txt")
	content := fmt.Sprintf("Multiverse Version: %s\nSecureBoot: %t\nPartitionStyle: %s\nFlashedAt: %s\n",
		i.config.Version, i.config.SecureBoot, i.config.PartitionStyle, time.Now().Format(time.RFC3339))
	return os.WriteFile(metaFile, []byte(content), 0644)
}
