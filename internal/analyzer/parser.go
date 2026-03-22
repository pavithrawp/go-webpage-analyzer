package analyzer

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// PageData holds all extracted data from a parsed HTML page
type PageData struct {
	HTMLVersion  string
	Title        string
	Headings     map[string]int
	Links        []string
	HasLoginForm bool
}

// parseHTML parses the HTML document and extracts all page data
func parseHTML(body string) (*PageData, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	data := &PageData{
		Headings: make(map[string]int),
	}

	// walk the HTML tree and extract data
	walkNode(doc, data)

	return data, nil
}

// walkNode recursively walks the HTML node tree and extracts page data
func walkNode(n *html.Node, data *PageData) {
	switch n.Type {
	case html.DoctypeNode:
		data.HTMLVersion = detectHTMLVersion(n)

	case html.ElementNode:
		switch n.Data {
		case "title":
			data.Title = extractText(n)
		case "h1", "h2", "h3", "h4", "h5", "h6":
			data.Headings[n.Data]++
		case "a":
			if href := getAttr(n, "href"); href != "" {
				data.Links = append(data.Links, href)
			}
		case "form":
			if hasPasswordInput(n) {
				data.HasLoginForm = true
			}
		}
	}

	// recursively visit all child nodes
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		walkNode(child, data)
	}
}

// extractText walks all child nodes and concatenates their text content
func extractText(n *html.Node) string {
	var text string
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			text += child.Data
		}
	}
	return strings.TrimSpace(text)
}

// detectHTMLVersion detects the HTML version from the DOCTYPE node
func detectHTMLVersion(n *html.Node) string {
	if len(n.Attr) == 0 {
		return "HTML5"
	}

	// Older HTML versions has public as a attribute
	publicAttr := ""
	for _, attr := range n.Attr {
		if attr.Key == "public" {
			publicAttr = strings.ToLower(attr.Val)
			break
		}
	}

	if publicAttr == "" {
		return "HTML5"
	}

	versions := []struct {
		keyword string
		result  string
	}{
		{"html 4.01", "HTML 4.01"},
		{"html 4.0", "HTML 4.0"},
		{"xhtml 1.1", "XHTML 1.1"},
		{"xhtml 1.0", "XHTML 1.0"},
		{"html 3.2", "HTML 3.2"},
		{"html 2.0", "HTML 2.0"},
	}

	variants := []struct {
		keyword string
		suffix  string
	}{
		{"strict", " Strict"},
		{"transitional", " Transitional"},
		{"frameset", " Frameset"},
	}

	for _, v := range versions {
		if strings.Contains(publicAttr, v.keyword) {
			for _, variant := range variants {
				if strings.Contains(publicAttr, variant.keyword) {
					return v.result + variant.suffix
				}
			}
			return v.result
		}
	}

	return "HTML5"
}

// getAttr returns the value of the given attribute from the node
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// hasPasswordInput checks if the given form node contains a password input
func hasPasswordInput(n *html.Node) bool {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "input" {
			if getAttr(child, "type") == "password" {
				return true
			}
		}
		if hasPasswordInput(child) {
			return true
		}
	}
	return false
}
