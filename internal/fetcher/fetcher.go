package fetcher

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/GangaRamPrasad2004/WallpaperAutomation/internal/config"
)

func GetDailyImageURL(cfg *config.Config) (string, error) {

	if cfg == nil {
		return "", fmt.Errorf("configuration cannot be nil")
	}

	if err := cfg.Validate(); err != nil {
		return "", fmt.Errorf("invalid configuration: %w", err)
	}

	today := time.Now().Format("2006-01-02")

	baseURL := strings.TrimRight(cfg.ImageSourceURL, "/")

	imageURL := fmt.Sprintf(
		"%s/%s/%d/%d",
		baseURL,
		today,
		cfg.ImageWidth,
		cfg.ImageHeight,
	)

	parsedURL, err := url.ParseRequestURI(imageURL)

	if err != nil {
		return "", fmt.Errorf("failed to generate image URL: %w", err)
	}

	return parsedURL.String(), nil
}
