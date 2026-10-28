# Multiverse Tools Suite

A Go-based desktop toolkit for creating bootable Multiverse drives, managing disk partitions, configuring plugins, and generating GRUB themes.

## Features

- **Disk enumeration** — list available storage drives as JSON.
- **Installer engine** — formats a target volume and flashes Multiverse boot structure (MBR/GPT, GRUB config, Adwaita theme, metadata).
- **Plugson** — native web-based plugin configurator UI for `ventoy.json` (themes, auto-install, persistence, passwords).
- **Theme maker** — generates a custom GRUB Adwaita theme.
- **Gio GUI** — native desktop GUI frontend.

## Commands

```
multiverse list                  # List available storage drives
multiverse install [device]      # Install Multiverse to target drive
multiverse plugson               # Launch native Plugson configurator UI
multiverse theme                # Generate custom GRUB Adwaita theme
multiverse gui                  # Launch Gio native GUI desktop app
multiverse version              # Print version information
```

## Build

Requires Go 1.24+.

```powershell
go build -o multiverse.exe ./cmd/multiverse
go test ./...
```

## Project layout

```
cmd/multiverse/     # CLI entrypoint
pkg/disk/           # Disk enumeration
pkg/installer/      # Installation engine
pkg/plugson/        # Plugin configurator (native + web)
pkg/theme/          # GRUB theme generator
pkg/gui/            # Gio desktop GUI
www/                # Plugson web frontend assets
multiverse_config/  # Sample ventoy.json
```

## License

See the Multiverse project for licensing details.
