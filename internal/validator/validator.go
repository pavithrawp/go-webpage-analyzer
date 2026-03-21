package validator

import (
	"fmt"
	"net/url"
	"regexp"
)

// the regex pattern for validating URLs
var urlPattern = regexp.MustCompile(`^https?://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=%]+$`)

// ValidateURL validates the given URL and returns an error if it is invalid
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL is required")
	}

	// check using regex first
	if !urlPattern.MatchString(rawURL) {
		return fmt.Errorf("invalid URL format: must start with http:// or https://")
	}

	// Regex provides fast format validation,
	// while net/url handles edge cases
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Host == "" {
		return fmt.Errorf("invalid URL: missing host")
	}

	return nil
}
