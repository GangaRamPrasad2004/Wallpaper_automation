package main

import (
	"fmt"

	"os"

	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/config"
	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/downloader"
	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/fetcher"
	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/logger"
	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/storage"
	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/validator"
	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/wallpaper"
)

func main() {

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Application failed:", err)
		os.Exit(1)
	}
}

func run() (runErr error) {

	// Load configuration
	cfg, err := config.LoadConfig()

	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Initialize logger
	appLogger, logFile, err := logger.NewLogger(cfg.LogDir)

	if err != nil {
		return fmt.Errorf("logger initialization failed: %w", err)
	}

	// Close log file after application execution
	defer logFile.Close()

	defer func() {
		if runErr != nil {
			appLogger.Error(
				"Wallpaper automation failed",
				"error", runErr,
			)
		}
	}()

	appLogger.Info("Wallpaper automation started")

	// Phase 3: Fetch image URL
	imageURL, err := fetcher.GetDailyImageURL(cfg)

	if err != nil {
		return fmt.Errorf("image fetching failed: %w", err)
	}

	appLogger.Info("Image URL generated", "url", imageURL)

	// Phase 4: Download image
	imagePath, err := downloader.DownloadImage(cfg, imageURL)

	if err != nil {
		return fmt.Errorf("image download failed: %w", err)
	}

	appLogger.Info("Image downloaded", "path", imagePath)

	// Phase 5: Validate image
	imageInfo, err := validator.ValidateImage(imagePath)

	if err != nil {
		return fmt.Errorf("image validation failed: %w", err)
	}

	appLogger.Info(
		"Image validation successful",
		"width", imageInfo.Width,
		"height", imageInfo.Height,
		"format", imageInfo.Format,
		"size_bytes", imageInfo.Size,
	)

	// Phase 6: Change wallpaper
	if err := wallpaper.ChangeWallpaper(imagePath); err != nil {
		return fmt.Errorf("wallpaper change failed: %w", err)
	}

	appLogger.Info(
		"Wallpaper changed successfully",
		"path", imagePath,
	)

	if err := storage.CleanupOldWallpapers(
		cfg.DownloadDir,
		imagePath,
	); err != nil {

		// Cleanup failure should not invalidate a successful
		// wallpaper update.
		appLogger.Warn(
			"Wallpaper cleanup encountered errors",
			"error", err,
		)
	} else {
		appLogger.Info(
			"Previous wallpapers moved to Recycle Bin",
		)
	}

	appLogger.Info("Wallpaper automation completed")

	return nil
}
