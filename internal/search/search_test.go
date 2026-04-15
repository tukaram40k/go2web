package search_test

import (
	"strings"
	"testing"

	"go2web/internal/search"
)

func TestBuildURL(t *testing.T) {
	got, err := search.BuildURL("golang http parser")
	if err != nil {
		t.Fatalf("BuildURL returned unexpected error: %v", err)
	}

	if !strings.HasPrefix(got, "https://duckduckgo.com/html/?q=") {
		t.Fatalf("unexpected base URL: %s", got)
	}

	if !strings.Contains(got, "golang+http+parser") {
		t.Fatalf("query term was not escaped as expected: %s", got)
	}
}

func TestBuildURLEmptyTerm(t *testing.T) {
	_, err := search.BuildURL("   ")
	if err == nil {
		t.Fatal("expected error for empty search term")
	}
}

func TestExtractResults(t *testing.T) {
	body := []byte(`
<html>
  <body>
    <a class="result__a" href="https://example.com/one">First Result</a>
    <a class="result__a" href="https://duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.org%2Ftwo">Second <b>Result</b></a>
    <a class="result__a" href="https://example.com/one">Duplicate URL</a>
  </body>
</html>`)

	results, err := search.ExtractResults(body, 10)
	if err != nil {
		t.Fatalf("ExtractResults returned unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 unique results, got %d", len(results))
	}

	if results[0].Title != "First Result" {
		t.Fatalf("unexpected first title: %q", results[0].Title)
	}

	if results[0].URL != "https://example.com/one" {
		t.Fatalf("unexpected first URL: %q", results[0].URL)
	}

	if results[1].Title != "Second Result" {
		t.Fatalf("unexpected second title: %q", results[1].Title)
	}

	if results[1].URL != "https://example.org/two" {
		t.Fatalf("unexpected second URL: %q", results[1].URL)
	}
}

func TestExtractResultsLimit(t *testing.T) {
	body := []byte(`
<a class="result__a" href="https://one.example">One</a>
<a class="result__a" href="https://two.example">Two</a>
<a class="result__a" href="https://three.example">Three</a>`)

	results, err := search.ExtractResults(body, 2)
	if err != nil {
		t.Fatalf("ExtractResults returned unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results due to limit, got %d", len(results))
	}
}
