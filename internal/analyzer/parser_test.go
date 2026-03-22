package analyzer

import (
	"testing"
	"strings"
)

// TestParseHTML_HTML5 tests that HTML5 is detected correctly
func TestParseHTML_HTML5(t *testing.T) {
	body := `<!DOCTYPE html>
	<html>
		<head><title>Test Page</title></head>
		<body><h1>Hello</h1></body>
	</html>`

	data, err := parseHTML(strings.NewReader(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data.HTMLVersion != "HTML5" {
		t.Errorf("expected HTML5, got %s", data.HTMLVersion)
	}
}

// TestParseHTML_HTML401 tests that HTML 4.01 is detected correctly
func TestParseHTML_HTML401(t *testing.T) {
	body := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"
	"http://www.w3.org/TR/html4/strict.dtd">
	<html>
		<head><title>Test Page</title></head>
		<body></body>
	</html>`

	data, err := parseHTML(strings.NewReader(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data.HTMLVersion != "HTML 4.01" {
		t.Errorf("expected HTML 4.01, got %s", data.HTMLVersion)
	}
}

// TestParseHTML_Title tests that the page title is extracted correctly
func TestParseHTML_Title(t *testing.T) {
	body := `<!DOCTYPE html>
	<html>
		<head><title>My Test Page</title></head>
		<body></body>
	</html>`

	data, err := parseHTML(strings.NewReader(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data.Title != "My Test Page" {
		t.Errorf("expected 'My Test Page', got %s", data.Title)
	}
}

// TestParseHTML_Headings tests that headings are counted correctly
func TestParseHTML_Headings(t *testing.T) {
	body := `<!DOCTYPE html>
	<html>
		<body>
			<h1>Heading 1</h1>
			<h1>Heading 1 again</h1>
			<h2>Heading 2</h2>
			<h3>Heading 3</h3>
		</body>
	</html>`

	data, err := parseHTML(strings.NewReader(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data.Headings["h1"] != 2 {
		t.Errorf("expected 2 h1 headings, got %d", data.Headings["h1"])
	}
	if data.Headings["h2"] != 1 {
		t.Errorf("expected 1 h2 heading, got %d", data.Headings["h2"])
	}
	if data.Headings["h3"] != 1 {
		t.Errorf("expected 1 h3 heading, got %d", data.Headings["h3"])
	}
}

// TestParseHTML_Links tests that links are extracted correctly
func TestParseHTML_Links(t *testing.T) {
	body := `<!DOCTYPE html>
	<html>
		<body>
			<a href="https://external.com">External</a>
			<a href="/internal">Internal</a>
			<a href="">Empty</a>
		</body>
	</html>`

	data, err := parseHTML(strings.NewReader(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// empty href should not be included
	if len(data.Links) != 2 {
		t.Errorf("expected 2 links, got %d", len(data.Links))
	}
}

// TestParseHTML_LoginForm tests that login forms are detected correctly
func TestParseHTML_LoginForm(t *testing.T) {
	body := `<!DOCTYPE html>
	<html>
		<body>
			<form>
				<input type="text" name="username"/>
				<input type="password" name="password"/>
				<button type="submit">Login</button>
			</form>
		</body>
	</html>`

	data, err := parseHTML(strings.NewReader(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !data.HasLoginForm {
		t.Error("expected login form to be detected")
	}
}

// TestParseHTML_NoLoginForm tests that pages without login forms return false
func TestParseHTML_NoLoginForm(t *testing.T) {
	body := `<!DOCTYPE html>
	<html>
		<body>
			<form>
				<input type="text" name="search"/>
				<button type="submit">Search</button>
			</form>
		</body>
	</html>`

	data, err := parseHTML(strings.NewReader(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data.HasLoginForm {
		t.Error("expected no login form to be detected")
	}
}
