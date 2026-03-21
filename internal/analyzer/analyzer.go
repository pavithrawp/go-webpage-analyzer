package analyzer

import (
	"fmt"
	"io"
)

// Result holds the full analysis result of a web page
type Result struct {
	HTMLVersion  string         `json:"html_version"`
	Title        string         `json:"title"`
	Headings     map[string]int `json:"headings"`
	Links        []string       `json:"links"`
	HasLoginForm bool           `json:"has_login_form"`
}

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
}

// Analyze fetches and analyzes the given URL
// TODO: add parser, link checker and login form detection
func (a *Analyzer) Analyze(url string) (*Result, error) {
	resp, err := fetchURL(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pageData, err := parseHTML(string(body), url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return &Result{
		HTMLVersion:  pageData.HTMLVersion,
		Title:        pageData.Title,
		Headings:     pageData.Headings,
		Links:        pageData.Links,
		HasLoginForm: pageData.HasLoginForm,
	}, nil
}
