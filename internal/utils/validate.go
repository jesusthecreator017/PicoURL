package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func ValidateURLFormat(rawUrl string) error {
	u, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("URL must include scheme and host")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be HTTP or HTTPS")
	}

	return nil
}

func ValidateURLReachable(rawURL string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(rawURL)
	if err != nil {
		return fmt.Errorf("URL is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("URL returned status: %d", resp.StatusCode)
	}
	return nil
}
