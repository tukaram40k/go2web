package cli_test

import (
	"flag"
	"testing"

	"go2web/internal/cli"
)

func TestGetFlags_Default(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	flags := cli.GetFlagsFromSet(fs, []string{})
	if flags == nil {
		t.Fatal("Expected flags to be non-nil")
	}

	if flags.URL != "http://example.com" {
		t.Errorf("Expected default URL 'http://example.com', got %q", flags.URL)
	}
	if flags.Search != "search term" {
		t.Errorf("Expected default Search 'search term', got %q", flags.Search)
	}
	if flags.Help != false {
		t.Errorf("Expected default Help false, got %v", flags.Help)
	}
}

func TestGetFlags_URL(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"custom URL with --url", []string{"-url", "https://google.com"}, "https://google.com"},
		{"custom URL with -u", []string{"-u", "https://google.com"}, "https://google.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			flags := cli.GetFlagsFromSet(fs, tt.args)
			if flags.URL != tt.expected {
				t.Errorf("Expected URL %q, got %q", tt.expected, flags.URL)
			}
		})
	}
}

func TestGetFlags_Search(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"custom search with --search", []string{"-search", "golang"}, "golang"},
		{"custom search with -s", []string{"-s", "golang"}, "golang"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			flags := cli.GetFlagsFromSet(fs, tt.args)
			if flags.Search != tt.expected {
				t.Errorf("Expected Search %q, got %q", tt.expected, flags.Search)
			}
		})
	}
}

func TestGetFlags_Help(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"help with --help", []string{"-help"}, true},
		{"help with -h", []string{"-h"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			flags := cli.GetFlagsFromSet(fs, tt.args)
			if flags.Help != tt.expected {
				t.Errorf("Expected Help %v, got %v", tt.expected, flags.Help)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		shouldError bool
	}{
		// Valid URLs
		{"valid https URL", "https://example.com", "https://example.com", false},
		{"valid http URL", "http://example.com", "http://example.com", false},
		{"valid URL with path", "https://example.com/path", "https://example.com/path", false},
		{"valid URL with query params", "https://example.com/search?q=golang", "https://example.com/search?q=golang", false},
		{"valid URL with subdomain", "https://www.example.com", "https://www.example.com", false},
		{"valid URL with multiple subdomains", "https://a.b.c.example.com", "https://a.b.c.example.com", false},

		// URLs that get http:// prepended
		{"URL without scheme", "example.com", "http://example.com", false},
		{"URL without scheme and path", "example.com/path", "http://example.com/path", false},

		// Invalid URLs
		{"empty URL", "", "", true},
		{"URL without domain", "http://", "", true},
		{"URL with invalid domain (single part)", "http://localhost", "", true},
		{"URL with trailing dot", "http://example.", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cli.ValidateURL(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for input %q, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Expected result %q, got %q", tt.expected, result)
				}
			}
		})
	}
}
