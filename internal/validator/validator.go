package validator

import (
	"fmt"
	"net/url"
)

// ValidateURL validates the given URL and returns an error if it is invalid
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL is required")
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid URL: must start with http:// or https://")
	}

	if parsed.Host == "" {
		return fmt.Errorf("invalid URL: missing host")
	}

	return nil
}
