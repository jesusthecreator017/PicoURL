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

const userAgent = "Mozilla/5.0 (compatible; PicoURL/1.0)"

func ValidateURLReachable(rawURL string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Try HEAD first — lightweight check
	headReq, err := http.NewRequest(http.MethodHead, rawURL, nil)
	if err != nil {
		return fmt.Errorf("URL is not reachable: %w", err)
	}
	headReq.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(headReq)
	if err != nil {
		// Network-level failure (DNS, connection refused, timeout) — truly unreachable
		return fmt.Errorf("URL is not reachable: %w", err)
	}
	resp.Body.Close()

	// Any successful response means the server is alive
	if resp.StatusCode < 400 {
		return nil
	}

	// Server responded with 4xx/5xx — some sites reject HEAD or block bots.
	// Fall back to GET to give it a second chance.
	getReq, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("URL is not reachable: %w", err)
	}
	getReq.Header.Set("User-Agent", userAgent)

	resp, err = client.Do(getReq)
	if err != nil {
		// Network failure on GET — truly unreachable
		return fmt.Errorf("URL is not reachable: %w", err)
	}
	resp.Body.Close()

	// If we got any HTTP response at all, the server is alive — the URL exists
	// behind access controls or bot protection, but it's not "unreachable".
	return nil
}
