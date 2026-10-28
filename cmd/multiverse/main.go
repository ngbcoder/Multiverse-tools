package main

import (
	"fmt"
	"os"

	"github.com/multiverse/suite/pkg/disk"
	"github.com/multiverse/suite/pkg/gui"
	"github.com/multiverse/suite/pkg/installer"
	"github.com/multiverse/suite/pkg/plugson"
	"github.com/multiverse/suite/pkg/theme"
)

const AppVersion = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "list":
		jsonStr, err := disk.FormatDisksJSON()
		if err != nil {
			fmt.Printf("Error listing disks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(jsonStr)

	case "install":
		dev := "/dev/sdb"
		if len(os.Args) > 2 {
			dev = os.Args[2]
		}
		inst := installer.NewInstaller(installer.Config{
			TargetDevice:   dev,
			SecureBoot:     true,
			PartitionStyle: "GPT",
			Version:        AppVersion,
		})
		fmt.Printf("Starting installation on %s...\n", dev)
		err := inst.Install(func(step string, percent int) {
			fmt.Printf("[%d%%] %s\n", percent, step)
		})
		if err != nil {
			fmt.Printf("Installation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Installation completed successfully!")

	case "plugson":
		if err := plugson.RunNativeUI(); err != nil {
			fmt.Printf("Plugson Native UI Error: %v\n", err)
			os.Exit(1)
		}

	case "theme":
		tm := theme.NewThemeMaker(theme.ThemeConfig{
			Title:        "Multiverse Custom Theme",
			PrimaryColor: "#2563eb",
		})
		outDir := "./theme_output"
		if err := tm.SaveTheme(outDir); err != nil {
			fmt.Printf("Failed to save theme: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Theme generated successfully in %s!\n", outDir)

	case "gui":
		if err := gui.RunGUI(); err != nil {
			fmt.Printf("GUI Error: %v\n", err)
			os.Exit(1)
		}

	case "version":
		fmt.Printf("Multiverse Tools Suite v%s (Golang)\n", AppVersion)

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Multiverse Desktop App Suite (Go)")
	fmt.Println("Usage: multiverse <command> [args]")
	fmt.Println("\nCommands:")
	fmt.Println("  list      - List available storage drives")
	fmt.Println("  install   - Install Multiverse to target drive [device]")
	fmt.Println("  plugson   - Launch native Plugson configurator UI")
	fmt.Println("  theme     - Generate custom GRUB Adwaita theme")
	fmt.Println("  gui       - Launch Gio native GUI desktop app")
	fmt.Println("  version   - Print version information")
}
