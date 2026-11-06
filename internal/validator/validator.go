package validator

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxImageSize   = 25 * 1024 * 1024
	maxImagePixels = 40_000_000
	minWidth       = 1280
	minHeight      = 720
)

type ImageInfo struct {
	Width  int
	Height int
	Format string
	Size   int64
}

func ValidateImage(imagePath string) (*ImageInfo, error) {

	// Check whether file exists
	fileInfo, err := os.Stat(imagePath)

	if err != nil {
		return nil, fmt.Errorf("image file not found: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("image path is not a regular file")
	}

	// Validate file size
	if fileInfo.Size() <= 0 {
		return nil, fmt.Errorf("image file is empty")
	}

	if fileInfo.Size() > maxImageSize {
		return nil, fmt.Errorf("image exceeds maximum size of 25 MB")
	}

	// Open image
	file, err := os.Open(imagePath)

	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}

	defer file.Close()

	// Read image header
	imgConfig, format, err := image.DecodeConfig(file)

	if err != nil {
		return nil, fmt.Errorf("invalid image header: %w", err)
	}

	// Validate supported format
	switch format {
	case "jpeg", "png", "gif":
		// Supported format
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	// Validate image dimensions
	width := imgConfig.Width
	height := imgConfig.Height

	if width < minWidth || height < minHeight {
		return nil, fmt.Errorf(
			"image resolution too low: %dx%d",
			width,
			height,
		)
	}

	// Prevent excessively large decoded images
	pixels := int64(width) * int64(height)

	if pixels > maxImagePixels {
		return nil, fmt.Errorf("image dimensions exceed maximum pixel limit")
	}

	// Validate file extension against actual image format
	expectedExtension := map[string]string{
		"jpeg": ".jpg",
		"png":  ".png",
		"gif":  ".gif",
	}[format]

	actualExtension := strings.ToLower(filepath.Ext(imagePath))

	if actualExtension != expectedExtension &&
		!(format == "jpeg" && actualExtension == ".jpeg") {
		return nil, fmt.Errorf(
			"image extension mismatch: expected %s, got %s",
			expectedExtension,
			actualExtension,
		)
	}

	// Reset file position
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	// Decode entire image to detect corrupted image data
	_, decodedFormat, err := image.Decode(file)

	if err != nil {
		return nil, fmt.Errorf("image decoding failed: %w", err)
	}

	if decodedFormat != format {
		return nil, fmt.Errorf("image format mismatch during decoding")
	}

	return &ImageInfo{
		Width:  width,
		Height: height,
		Format: format,
		Size:   fileInfo.Size(),
	}, nil
}
