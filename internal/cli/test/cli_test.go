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

	if flags.URL != "" {
		t.Errorf("Expected empty URL, got %q", flags.URL)
	}
	if flags.Search != "" {
		t.Errorf("Expected empty Search, got %q", flags.Search)
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
		{"Invalid URL 4", "http://example.com.http://example.com", "", true},
		{"Invalid URL 5", "http://example.", "", true},
		{"Invalid URL 6", "http://example.com.", "", true},
		{"Invalid URL 7", "http://.com", "", true},
		{"Invalid URL 8", "http://example.com:99999", "", true},
		{"Invalid URL 9", "https://example..com", "", true},
		{"Invalid URL 10", "http://-example.com", "", true},
		{"Invalid URL 11", "javascript:alert(1)", "", true},
		{"Invalid URL 12", "http://[::1]]", "", true},
		{"Invalid URL 13", "http://123.456.789.0", "", true},
		{"Invalid URL 14", "http:// user:pass@example.com", "", true},
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
