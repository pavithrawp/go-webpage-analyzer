package analyzer

import (
	"fmt"
	"net/http"
	"time"
)

// Default timeout for the URL to respond
const defaultTimeout = 10 * time.Second

// fetchURL fetches the HTML content from the given URL and throw an error if the URL is unreachable
func fetchURL(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: defaultTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("URL returned status %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return resp, nil
}
