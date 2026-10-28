package installer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallerFormatAndFlashing(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "multiverse_format_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a dummy file to ensure formatting cleans it up
	dummyFile := filepath.Join(tmpDir, "old_file.txt")
	_ = os.WriteFile(dummyFile, []byte("old data"), 0644)

	cfg := Config{
		TargetDevice:   tmpDir,
		SecureBoot:     true,
		PartitionStyle: "GPT",
		Version:        "v1.0.0",
		FileSystem:     "exFAT",
	}

	inst := NewInstaller(cfg)

	var lastPercent int
	err = inst.Install(func(step string, percent int) {
		lastPercent = percent
	})
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}
	if lastPercent != 100 {
		t.Errorf("expected final progress 100, got %d", lastPercent)
	}

	// Verify old file was cleared by format step
	if _, err := os.Stat(dummyFile); !os.IsNotExist(err) {
		t.Errorf("dummy file was not removed during formatting")
	}

	// Verify boot sector file written
	bootBin := filepath.Join(tmpDir, "multiverse_boot.bin")
	if _, err := os.Stat(bootBin); os.IsNotExist(err) {
		t.Errorf("boot sector signature file not found")
	}

	// Verify grub.cfg written
	grubCfg := filepath.Join(tmpDir, "boot", "grub", "grub.cfg")
	if _, err := os.Stat(grubCfg); os.IsNotExist(err) {
		t.Errorf("grub.cfg was not written to disk")
	}

	// Verify metadata file written
	metaPath := filepath.Join(tmpDir, "multiverse_version.txt")
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		t.Errorf("metadata file was not created")
	}
}
