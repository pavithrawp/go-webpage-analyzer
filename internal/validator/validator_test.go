package validator

import (
	"testing"
)

// TestValidateURL tests that valid URLs pass validation
func TestValidateURL_Valid(t *testing.T) {
	validURLs := []string{
		"https://example.com",
		"http://example.com",
		"https://example.com/path",
		"https://abc.com/path?query=1",
		"https://sub.example.com",
	}

	for _, url := range validURLs {
		if err := ValidateURL(url); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", url, err)
		}
	}
}

// TestValidateURL_Empty tests that empty URLs fail validation
func TestValidateURL_Empty(t *testing.T) {
	if err := ValidateURL(""); err == nil {
		t.Error("expected error for empty URL, got nil")
	}
}

// TestValidateURL_InvalidFormat tests that invalid URLs fail validation
func TestValidateURL_InvalidFormat(t *testing.T) {
	invalidURLs := []string{
		"dawsda2-2",
		"ftp://example.com",
		"alert(1)",
		"http://",
		"//abc.com",
	}

	for _, url := range invalidURLs {
		if err := ValidateURL(url); err == nil {
			t.Errorf("expected %q to be invalid, got nil", url)
		}
	}
}
