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
		t.Fatalf("expected 2 headers, got %d", len(resp.HeaderFields))
	}

	if resp.ContentType != "text/html" {
		t.Fatalf("unexpected content type: %q", resp.ContentType)
	}

	if !bytes.Equal(resp.Body, []byte("Hello")) {
		t.Fatalf("unexpected body: %q", string(resp.Body))
	}

	if !resp.ResponseIsOK {
		t.Fatal("expected ResponseIsOK to be true for 200 response")
	}

	if resp.IsRedirected {
		t.Fatal("expected IsRedirected to be false")
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

func TestParseChunkedBody(t *testing.T) {
	raw := []byte("HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\nContent-Type: text/plain\r\n\r\n5\r\nHello\r\n6\r\n World\r\n0\r\n\r\n")

	resp, err := parser.Parse(raw)
	if err != nil {
		t.Fatalf("Parse returned unexpected error: %v", err)
	}

	if got := string(resp.Body); got != "Hello World" {
		t.Fatalf("unexpected decoded chunked body: %q", got)
	}
}

func TestParseWithRedirectInfo(t *testing.T) {
	raw := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello")

	resp, err := parser.ParseWithRedirectInfo(raw, 2)
	if err != nil {
		t.Fatalf("ParseWithRedirectInfo returned unexpected error: %v", err)
	}

	if !resp.IsRedirected {
		t.Fatal("expected IsRedirected to be true")
	}

	if resp.RedirectCount != 2 {
		t.Fatalf("expected RedirectCount 2, got %d", resp.RedirectCount)
	}
}
