package analyzer

import (
	"fmt"
	"net/http"
	"time"
)

// FetchError represents an error that occurred while fetching a URL.
type FetchError struct {
	StatusCode int
	Message    string
}

// Default timeout for the URL to respond
const defaultTimeout = 10 * time.Second

// fetchURL fetches the HTML content from the given URL and throw an error if the URL is unreachable
// Note: this uses a standard HTTP client and does not execute JavaScript.
// Pages that rely on client-side rendering (React, Angular, Vue) may return
// incomplete content as the JavaScript is not executed before parsing.
func fetchURL(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: defaultTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, &FetchError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	return resp, nil
}

// Error implements the error interface
func (e *FetchError) Error() string {
	return fmt.Sprintf("URL returned status %d: %s", e.StatusCode, e.Message)
}
