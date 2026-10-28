# Multiverse Go Application Suite Walkthrough

We have added target volume formatting to the **Multiverse** installer engine.

## Automatic Volume Formatting (`pkg/installer/installer.go`)

- **Format Action (`FormatTargetVolume`)**:
  - Invokes native Windows formatting utility (`format.com <drive> /FS:exFAT /Q /Y /V:MULTIVERSE`) on the target volume.
  - Clears out all existing files/directories on the target drive to guarantee a clean file system state before deploying Multiverse files.
  - Supports `exFAT`, `FAT32`, and `NTFS` file system options.

## Verification Results

- **Build**:
  ```powershell
  go build -o multiverse.exe ./cmd/multiverse
  ```
- **Unit Tests**:
  ```powershell
  go test ./...
  ok  github.com/multiverse/suite/pkg/disk
  ok  github.com/multiverse/suite/pkg/installer
  ok  github.com/multiverse/suite/pkg/plugson
  ok  github.com/multiverse/suite/pkg/theme
  ```
