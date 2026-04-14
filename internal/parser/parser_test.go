package parser_test

import (
	"bytes"
	"testing"

	"go2web/internal/parser"
)

func TestParseCRLFResponse(t *testing.T) {
	raw := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 5\r\n\r\nHello")

	resp, err := parser.Parse(raw)
	if err != nil {
		t.Fatalf("Parse returned unexpected error: %v", err)
	}

	if resp.StatusLine != "HTTP/1.1 200 OK" {
		t.Fatalf("unexpected status line: %q", resp.StatusLine)
	}

	if len(resp.HeaderFields) != 2 {
		t.Fatalf("expected 2 header fields, got %d", len(resp.HeaderFields))
	}

	if resp.HeaderFields[0] != "Content-Type: text/html" {
		t.Fatalf("unexpected first header: %q", resp.HeaderFields[0])
	}

	if !bytes.Equal(resp.Body, []byte("Hello")) {
		t.Fatalf("unexpected body: %q", string(resp.Body))
	}

	if !resp.ResponseIsOK {
		t.Fatal("expected ResponseIsOK to be true for 200 response")
	}
}

func TestParseLFOnlyResponse(t *testing.T) {
	raw := []byte("HTTP/1.1 404 Not Found\nContent-Type: text/plain\n\nMissing")

	resp, err := parser.Parse(raw)
	if err != nil {
		t.Fatalf("Parse returned unexpected error: %v", err)
	}

	if resp.StatusLine != "HTTP/1.1 404 Not Found" {
		t.Fatalf("unexpected status line: %q", resp.StatusLine)
	}

	if len(resp.HeaderFields) != 1 {
		t.Fatalf("expected 1 header field, got %d", len(resp.HeaderFields))
	}

	if !bytes.Equal(resp.Body, []byte("Missing")) {
		t.Fatalf("unexpected body: %q", string(resp.Body))
	}

	if resp.ResponseIsOK {
		t.Fatal("expected ResponseIsOK to be false for 404 response")
	}
}

func TestParseMissingSeparator(t *testing.T) {
	raw := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain")

	_, err := parser.Parse(raw)
	if err == nil {
		t.Fatal("expected parse error for missing header separator")
	}
}
