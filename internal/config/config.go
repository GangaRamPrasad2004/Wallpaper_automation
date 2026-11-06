package config

import (
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	ImageSourceURL string
	ImageWidth     int
	ImageHeight    int
	DownloadDir    string
	LogDir         string
	HTTPTimeout    time.Duration
	UserAgent      string
}

func LoadConfig() (*Config, error) {

	cfg := &Config{
		ImageSourceURL: "https://www.iskconbangalore.org/daily-darshan/",
		ImageWidth:     1920,
		ImageHeight:    1080,
		DownloadDir:    "wallpapers",
		LogDir:         "logs",
		HTTPTimeout:    30 * time.Second,
		UserAgent:      "DailyWallpaperBot/1.0",
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {

	parsedURL, err := url.ParseRequestURI(c.ImageSourceURL)

	if err != nil {
		return fmt.Errorf("invalid image source URL: %w", err)
	}

	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return fmt.Errorf("image source must use HTTP or HTTPS")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("image source host is missing")
	}

	if c.ImageWidth <= 0 || c.ImageHeight <= 0 {
		return fmt.Errorf("image dimensions must be positive")
	}

	if c.DownloadDir == "" || c.LogDir == "" {
		return fmt.Errorf("download and log directories are required")
	}

	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTP timeout must be positive")
	}

	if c.UserAgent == "" {
		return fmt.Errorf("user agent cannot be empty")
	}

	return nil
}
