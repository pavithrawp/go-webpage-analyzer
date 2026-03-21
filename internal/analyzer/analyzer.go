package analyzer

import (
	"io"
)

// Result holds the full analysis result of a web page
type Result struct {
	RawHTML string
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

	return &Result{
		RawHTML: string(body),
	}, nil
}
