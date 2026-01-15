package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var ErrEmptyURL = errors.New("empty URL")

// normalizeURL adds https:// prefix if missing
func normalizeURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrEmptyURL
	}
	// Add https:// prefix if URL scheme is missing
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}
	_, err := url.Parse(raw)
	return raw, err
}

// CheckOnce checks a single URL and returns status code, status text, and error
func CheckOnce(ctx context.Context, client *http.Client, rawURL string) (int, string, error) {
	u, err := normalizeURL(rawURL)
	if err != nil {
		return 0, "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	return resp.StatusCode, resp.Status, nil
}
