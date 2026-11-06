//go:build windows

package wallpaper

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	SPI_SETDESKWALLPAPER = 20

	SPIF_UPDATEINIFILE = 1
	SPIF_SENDCHANGE    = 2
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	systemParametersInfo = user32.NewProc(
		"SystemParametersInfoW",
	)
)

func ChangeWallpaper(imagePath string) error {

	// Convert to absolute path
	absolutePath, err := filepath.Abs(imagePath)

	if err != nil {
		return fmt.Errorf("failed to resolve image path: %w", err)
	}

	// Check whether the image exists
	fileInfo, err := os.Stat(absolutePath)

	if err != nil {
		return fmt.Errorf("wallpaper image not found: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("wallpaper path is not a regular file")
	}

	// Convert image path into UTF-16
	pathPointer, err := syscall.UTF16PtrFromString(absolutePath)

	if err != nil {
		return fmt.Errorf("failed to convert image path: %w", err)
	}

	// Call Windows API
	result, _, callErr := systemParametersInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		0,
		uintptr(unsafe.Pointer(pathPointer)),
		uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
	)

	// Check Windows API result
	if result == 0 {
		if callErr != nil && callErr != syscall.Errno(0) {
			return fmt.Errorf("Windows API failed: %w", callErr)
		}

		return fmt.Errorf("Windows API failed to change wallpaper")
	}

	return nil
}
