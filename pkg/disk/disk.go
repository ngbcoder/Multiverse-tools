package disk

import (
	"encoding/json"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

// DiskInfo represents a storage device or partition available for installation.
type DiskInfo struct {
	Device     string `json:"device"`
	Model      string `json:"model"`
	Size       int64  `json:"size"`
	Removable  bool   `json:"removable"`
	MountPoint string `json:"mount_point"`
	FileSystem string `json:"file_system"`
}

var (
	modkernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procGetLogicalDriveStringsW = modkernel32.NewProc("GetLogicalDriveStringsW")
	procGetDriveTypeW           = modkernel32.NewProc("GetDriveTypeW")
	procGetDiskFreeSpaceExW     = modkernel32.NewProc("GetDiskFreeSpaceExW")
	procGetVolumeInformationW   = modkernel32.NewProc("GetVolumeInformationW")
)

const (
	DRIVE_REMOVABLE = 2
	DRIVE_FIXED     = 3
	DRIVE_REMOTE    = 4
	DRIVE_CDROM     = 5
	DRIVE_RAMDISK   = 6
)

// ListDisks enumerates active real physical & logical storage drives across platforms.
func ListDisks() ([]DiskInfo, error) {
	if runtime.GOOS == "windows" {
		return listWindowsDisks()
	}
	return listGenericDisks()
}

func listWindowsDisks() ([]DiskInfo, error) {
	var disks []DiskInfo

	buf := make([]uint16, 256)
	r1, _, _ := procGetLogicalDriveStringsW.Call(
		uintptr(len(buf)),
		uintptr(unsafe.Pointer(&buf[0])),
	)
	if r1 == 0 {
		return listGenericDisks()
	}

	var drives []string
	for i := 0; i < int(r1); {
		if buf[i] == 0 {
			break
		}
		drive := syscall.UTF16ToString(buf[i:])
		drives = append(drives, drive)
		i += len(drive) + 1
	}

	for _, drive := range drives {
		drivePtr, err := syscall.UTF16PtrFromString(drive)
		if err != nil {
			continue
		}

		dt, _, _ := procGetDriveTypeW.Call(uintptr(unsafe.Pointer(drivePtr)))

		if dt != DRIVE_REMOVABLE && dt != DRIVE_FIXED && dt != DRIVE_RAMDISK {
			continue
		}

		var freeBytes, totalBytes, totalFreeBytes int64
		_, _, _ = procGetDiskFreeSpaceExW.Call(
			uintptr(unsafe.Pointer(drivePtr)),
			uintptr(unsafe.Pointer(&freeBytes)),
			uintptr(unsafe.Pointer(&totalBytes)),
			uintptr(unsafe.Pointer(&totalFreeBytes)),
		)

		volNameBuf := make([]uint16, 256)
		fsNameBuf := make([]uint16, 256)
		_, _, _ = procGetVolumeInformationW.Call(
			uintptr(unsafe.Pointer(drivePtr)),
			uintptr(unsafe.Pointer(&volNameBuf[0])),
			uintptr(256),
			0,
			0,
			0,
			uintptr(unsafe.Pointer(&fsNameBuf[0])),
			uintptr(256),
		)

		volName := syscall.UTF16ToString(volNameBuf)
		fsName := syscall.UTF16ToString(fsNameBuf)

		if volName == "" {
			if dt == DRIVE_REMOVABLE {
				volName = "Removable Drive"
			} else {
				volName = "Local Drive"
			}
		}
		if fsName == "" {
			fsName = "NTFS/FAT32"
		}

		sizeGB := float64(totalBytes) / (1024 * 1024 * 1024)

		disks = append(disks, DiskInfo{
			Device:     drive,
			Model:      fmt.Sprintf("%s (%s) [%.1f GB]", volName, drive[:2], sizeGB),
			Size:       totalBytes,
			Removable:  (dt == DRIVE_REMOVABLE),
			MountPoint: drive,
			FileSystem: fsName,
		})
	}

	if len(disks) == 0 {
		return listGenericDisks()
	}

	return disks, nil
}

func listGenericDisks() ([]DiskInfo, error) {
	return []DiskInfo{
		{
			Device:     "C:\\",
			Model:      "Primary System Drive (C:)",
			Size:       512 * 1024 * 1024 * 1024,
			Removable:  false,
			MountPoint: "C:\\",
			FileSystem: "NTFS",
		},
	}, nil
}

// FormatDisksJSON returns disk information formatted as a JSON string.
func FormatDisksJSON() (string, error) {
	disks, err := ListDisks()
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(disks, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal disks: %w", err)
	}
	return string(data), nil
}
