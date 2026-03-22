package analyzer

import (
	"context"
	"fmt"
	"io"
)

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

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
}

// Analyze fetches and analyzes the given URL
func (a *Analyzer) Analyze(ctx context.Context, url string) (*Result, error) {
	resp, err := fetchURL(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	pageData, err := parseHTML(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// concurrently check all links
	linkSummary := checkLinks(ctx, pageData.Links, url)

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
