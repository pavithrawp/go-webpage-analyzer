package analyzer

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// the number of concurrent workers for link checking
const workerCount = 20

// holds the result of checking a single link
type linkResult struct {
	URL          string
	IsInternal   bool
	IsAccessible bool
}

// holds the final summary of all links
type linkSummary struct {
	InternalCount     int
	ExternalCount     int
	InaccessibleCount int
	Results           []linkResult
}

// checkLinks classifies and concurrently checks all links for accessibility
func (a *Analyzer) checkLinks(ctx context.Context, links []string, baseURL string) *linkSummary {
	if len(links) == 0 {
		return &linkSummary{}
	}

	summary := &linkSummary{}
	seen := make(map[string]struct{})
	var uniqueLinks []string
	for _, link := range links {
		resolved := resolveLink(link, baseURL)
		if resolved == "" {
			continue
		}

		// count ALL links including duplicates
		if isInternalLink(resolved, baseURL) {
			summary.InternalCount++
		} else {
			summary.ExternalCount++
		}

		// deduplicate only for accessibility checking
		if _, exists := seen[resolved]; !exists {
			seen[resolved] = struct{}{}
			uniqueLinks = append(uniqueLinks, resolved)
		}
	}

	if len(uniqueLinks) == 0 {
		return summary
	}

	// check accessibility only for unique links
	// set the buffer size to the number of links. so that we can send all links into the channel without blocking
	jobs := make(chan string, len(uniqueLinks))
	results := make(chan linkResult, len(uniqueLinks))

	// start worker pool
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for link := range jobs {
				isAccessible := a.isLinkAccessible(ctx, link)
				results <- linkResult{
					URL:          link,
					IsAccessible: isAccessible,
				}
			}
		}()
	}

	// send resolved links to workers
	for _, link := range uniqueLinks {
		jobs <- link
	}
	close(jobs)

	// wait for all workers to finish then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if !result.IsAccessible {
			summary.InaccessibleCount++
		}
	}

	return summary
}

// isInternalLink checks if the link belongs to the same domain as the base URL
func isInternalLink(link, baseURL string) bool {
	base, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	parsed, err := url.Parse(link)
	if err != nil {
		return false
	}

	// relative URLs are internal
	if !parsed.IsAbs() {
		return true
	}

	return parsed.Host == base.Host
}

// isLinkAccessible checks if the link is accessible by making a HEAD request.
// Falls back to GET if HEAD returns 405 Method Not Allowed.
func (a *Analyzer) isLinkAccessible(ctx context.Context, link string) bool {
	accessible, status := a.checkStatus(ctx, http.MethodHead, link)
	if !accessible && status == http.StatusMethodNotAllowed {
		accessible, _ = a.checkStatus(ctx, http.MethodGet, link)
	}
	return accessible
}

// checkStatus makes an HTTP request and returns whether it succeeded and the status code
func (a *Analyzer) checkStatus(ctx context.Context, method, link string) (bool, int) {
	req, err := http.NewRequestWithContext(ctx, method, link, nil)
	if err != nil {
		return false, 0
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return false, 0
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode < 400, resp.StatusCode
}

// resolveLink converts relative URLs to full URLs
func resolveLink(link, baseURL string) string {

	// ignore anchor links
	if len(link) == 0 || link[0] == '#' {
		return ""
	}

	// filter out non HTTP schemes
	lower := strings.ToLower(link)
	if strings.HasPrefix(lower, "mailto:") ||
		strings.HasPrefix(lower, "javascript:") ||
		strings.HasPrefix(lower, "tel:") ||
		strings.HasPrefix(lower, "data:") ||
		strings.HasPrefix(lower, "ftp:") {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	parsed, err := url.Parse(link)
	if err != nil {
		return ""
	}

	// resolve relative URLs against base
	return base.ResolveReference(parsed).String()
}
