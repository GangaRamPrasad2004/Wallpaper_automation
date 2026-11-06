package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func NewLogger(logDir string) (*slog.Logger, *os.File, error) {

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, nil, fmt.Errorf(
			"failed to create log directory: %w",
			err,
		)
	}

	// Open log file in append mode
	logPath := filepath.Join(logDir, "wallpaper.log")

	logFile, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to open log file: %w",
			err,
		)
	}

	// Write logs to both terminal and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Configure structured logger
	handler := slog.NewTextHandler(
		multiWriter,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	appLogger := slog.New(handler)

	return appLogger, logFile, nil
}
