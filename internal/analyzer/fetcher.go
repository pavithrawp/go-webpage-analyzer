package analyzer

import (
	"context"
	"fmt"
	"net/http"
)

// FetchError represents an error that occurred while fetching a URL.
type FetchError struct {
	StatusCode int
	Message    string
}

// Error implements the error interface
func (e *FetchError) Error() string {
	return fmt.Sprintf("URL returned status %d: %s", e.StatusCode, e.Message)
}

// fetchURL fetches the HTML content from the given URL and throw an error if the URL is unreachable
// Note: this uses a standard HTTP client and does not execute JavaScript.
// Pages that rely on client-side rendering (React, Angular, Vue) may return
// incomplete content as the JavaScript is not executed before parsing.
func (a *Analyzer) fetchURL(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach URL: %w", err)
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
