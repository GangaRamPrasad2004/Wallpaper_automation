package downloader

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/config"
)

const maxImageSize = 25 * 1024 * 1024 // 25 MB

func DownloadImage(cfg *config.Config, imageURL string) (string, error) {

	if cfg == nil {
		return "", fmt.Errorf("configuration cannot be nil")
	}

	if err := cfg.Validate(); err != nil {
		return "", fmt.Errorf("invalid configuration: %w", err)
	}

	if strings.TrimSpace(imageURL) == "" {
		return "", fmt.Errorf("image URL cannot be empty")
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodGet, imageURL, nil)

	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", cfg.UserAgent)
	req.Header.Set("Accept", "image/*")

	// Create HTTP client
	client := &http.Client{
		Timeout: cfg.HTTPTimeout,
	}

	// Execute HTTP request
	resp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}

	defer resp.Body.Close()

	// Validate HTTP status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"unexpected HTTP status: %s",
			resp.Status,
		)
	}

	// Validate content type
	mediaType, _, err := mime.ParseMediaType(
		resp.Header.Get("Content-Type"),
	)

	if err != nil {
		return "", fmt.Errorf("invalid content type: %w", err)
	}

	extension, err := getImageExtension(mediaType)

	if err != nil {
		return "", err
	}

	// Check content length when available
	if resp.ContentLength > maxImageSize {
		return "", fmt.Errorf("image exceeds maximum size of 25 MB")
	}

	// Create download directory
	if err := os.MkdirAll(cfg.DownloadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create temporary file
	tempFile, err := os.CreateTemp(
		cfg.DownloadDir,
		".wallpaper-*.tmp",
	)

	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	tempPath := tempFile.Name()

	defer func() {
		tempFile.Close()
		os.Remove(tempPath)
	}()

	// Copy response body into temporary file
	size, err := io.Copy(
		tempFile,
		io.LimitReader(resp.Body, maxImageSize+1),
	)

	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	if size == 0 {
		return "", fmt.Errorf("downloaded image is empty")
	}

	if size > maxImageSize {
		return "", fmt.Errorf("image exceeds maximum size of 25 MB")
	}

	// Ensure data is flushed to disk
	if err := tempFile.Sync(); err != nil {
		return "", fmt.Errorf("failed to sync image file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close image file: %w", err)
	}

	// Generate unique filename
	filename := time.Now().Format("2006-01-02-150405.000000000") + extension

	finalPath := filepath.Join(cfg.DownloadDir, filename)

	// Move temporary file to final location
	if err := os.Rename(tempPath, finalPath); err != nil {
		return "", fmt.Errorf("failed to save final image: %w", err)
	}

	return finalPath, nil
}

func getImageExtension(mediaType string) (string, error) {

	switch mediaType {

	case "image/jpeg":
		return ".jpg", nil

	case "image/png":
		return ".png", nil

	case "image/gif":
		return ".gif", nil

	default:
		return "", fmt.Errorf(
			"unsupported image format: %s",
			mediaType,
		)
	}
}
