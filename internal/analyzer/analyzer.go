package analyzer

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// the maximum time allowed for a URL to respond
const defaultTimeout = 10 * time.Second

// Result holds the full analysis result of a web page
type Result struct {
	HTMLVersion       string         `json:"html_version"`
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	HasLoginForm      bool           `json:"has_login_form"`
	InternalLinks     int            `json:"internal_links"`
	ExternalLinks     int            `json:"external_links"`
	InaccessibleLinks int            `json:"inaccessible_links"`
}

type Analyzer struct {
	httpClient *http.Client
}

func New() *Analyzer {
	return &Analyzer{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// Analyze fetches and analyzes the given URL
func (a *Analyzer) Analyze(ctx context.Context, url string) (*Result, error) {
	resp, err := a.fetchURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	pageData, err := parseHTML(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// concurrently check all links
	linkSummary := a.checkLinks(ctx, pageData.Links, url)

	return &Result{
		HTMLVersion:       pageData.HTMLVersion,
		Title:             pageData.Title,
		Headings:          pageData.Headings,
		HasLoginForm:      pageData.HasLoginForm,
		InternalLinks:     linkSummary.InternalCount,
		ExternalLinks:     linkSummary.ExternalCount,
		InaccessibleLinks: linkSummary.InaccessibleCount,
	}, nil
}
