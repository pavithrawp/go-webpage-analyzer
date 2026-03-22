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
type LinkResult struct {
	URL          string
	IsInternal   bool
	IsAccessible bool
}

// holds the final summary of all links
type LinkSummary struct {
	InternalCount     int
	ExternalCount     int
	InaccessibleCount int
	Results           []LinkResult
}

// checkLinks classifies and concurrently checks all links for accessibility
func (a *Analyzer) checkLinks(ctx context.Context, links []string, baseURL string) *LinkSummary {
	if len(links) == 0 {
		return &LinkSummary{}
	}

	// converts relative URLs to full URLs
	var resolvedLinks []string
	for _, link := range links {
		resolved := resolveLink(link, baseURL)
		if resolved != "" {
			resolvedLinks = append(resolvedLinks, resolved)
		}
	}

	if len(resolvedLinks) == 0 {
		return &LinkSummary{}
	}

	// set the buffer size to the number of links. so that we can send all links into the channel without blocking
	jobs := make(chan string, len(resolvedLinks))
	results := make(chan LinkResult, len(resolvedLinks))

	// start worker pool
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for link := range jobs {
				isInternal := isInternalLink(link, baseURL)
				isAccessible := a.isLinkAccessible(ctx, link)
				results <- LinkResult{
					URL:          link,
					IsInternal:   isInternal,
					IsAccessible: isAccessible,
				}
			}
		}()
	}

	// send resolved links to workers
	for _, link := range resolvedLinks {
		jobs <- link
	}
	close(jobs)

	// wait for all workers to finish then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect results
	summary := &LinkSummary{}
	for result := range results {
		summary.Results = append(summary.Results, result)
		if result.IsInternal {
			summary.InternalCount++
		} else {
			summary.ExternalCount++
		}
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

// isLinkAccessible checks if the link is accessible by making a HEAD request
func (a *Analyzer) isLinkAccessible(ctx context.Context, link string) bool {
	// used head just know if the link is alive.
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, link, nil)
	if err != nil {
		return false
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode < 400
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
