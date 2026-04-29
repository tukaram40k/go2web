package parser

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

type Response struct {
	StatusLine      string
	HeaderFields    []string
	Body            []byte
	ResponseIsOK    bool
	IsRedirected    bool
	RedirectCount   int
	ContentType     string
	Cached          bool
	CacheAgeMinutes int
}

func parseStatusCode(statusLine string) (int, error) {
	parts := strings.Fields(statusLine)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid status line: %q", statusLine)
	}

	code, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid status code in status line: %w", err)
	}

	return code, nil
}

func splitHeadersAndBody(raw []byte) (string, []string, []byte, error) {
	headerEnd := bytes.Index(raw, []byte("\r\n\r\n"))
	lineSeparator := "\r\n"
	bodyStart := 4

	if headerEnd == -1 {
		headerEnd = bytes.Index(raw, []byte("\n\n"))
		lineSeparator = "\n"
		bodyStart = 2
	}

	if headerEnd == -1 {
		return "", nil, nil, fmt.Errorf("invalid http response: missing header separator")
	}

	lines := strings.Split(string(raw[:headerEnd]), lineSeparator)
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return "", nil, nil, fmt.Errorf("invalid http response: missing status line")
	}

	bodyIndex := headerEnd + bodyStart
	if bodyIndex > len(raw) {
		return "", nil, nil, fmt.Errorf("invalid http response: body index out of range")
	}

	statusLine := lines[0]
	headers := make([]string, 0, len(lines)-1)
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		headers = append(headers, line)
	}

	return statusLine, headers, raw[bodyIndex:], nil
}

func parseHeaders(lines []string) map[string][]string {
	headers := make(map[string][]string)

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		canonicalName := textproto.CanonicalMIMEHeaderKey(name)

		headers[canonicalName] = append(headers[canonicalName], value)
	}

	return headers
}

func hasTokenCI(values []string, token string) bool {
	for _, v := range values {
		parts := strings.Split(v, ",")
		for _, part := range parts {
			if strings.EqualFold(strings.TrimSpace(part), token) {
				return true
			}
		}
	}

	return false
}

func decodeChunkedBody(body []byte) ([]byte, error) {
	reader := bufio.NewReader(bytes.NewReader(body))
	var decoded bytes.Buffer

	for {
		sizeLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("invalid chunked body: failed to read chunk size: %w", err)
		}

		sizeLine = strings.TrimSpace(sizeLine)
		if sizeLine == "" {
			continue
		}

		if idx := strings.Index(sizeLine, ";"); idx != -1 {
			sizeLine = sizeLine[:idx]
		}

		chunkSize, err := strconv.ParseInt(sizeLine, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid chunk size %q: %w", sizeLine, err)
		}

		if chunkSize == 0 {
			_, _ = reader.ReadString('\n')
			break
		}

		if chunkSize < 0 {
			return nil, fmt.Errorf("invalid negative chunk size")
		}

		if _, err := io.CopyN(&decoded, reader, chunkSize); err != nil {
			return nil, fmt.Errorf("invalid chunked body: truncated chunk: %w", err)
		}

		lineEnd, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("invalid chunked body: missing chunk terminator: %w", err)
		}

		if lineEnd != "\r\n" && lineEnd != "\n" {
			return nil, fmt.Errorf("invalid chunked body: invalid chunk terminator %q", lineEnd)
		}
	}

	return decoded.Bytes(), nil
}

func decodeBodyByHeaders(headers map[string][]string, body []byte) ([]byte, error) {
	decodedBody := body

	if hasTokenCI(headers["Transfer-Encoding"], "chunked") {
		chunkedBody, err := decodeChunkedBody(body)
		if err != nil {
			return nil, err
		}

		decodedBody = chunkedBody
	} else if lengthValues, ok := headers["Content-Length"]; ok && len(lengthValues) > 0 {
		contentLength, err := strconv.Atoi(strings.TrimSpace(lengthValues[0]))
		if err == nil {
			if contentLength >= 0 && contentLength <= len(decodedBody) {
				decodedBody = decodedBody[:contentLength]
			}
		}
	}

	if hasTokenCI(headers["Content-Encoding"], "gzip") {
		gzReader, err := gzip.NewReader(bytes.NewReader(decodedBody))
		if err != nil {
			return nil, fmt.Errorf("failed to initialize gzip reader: %w", err)
		} else {
			unzipped, readErr := io.ReadAll(gzReader)
			closeErr := gzReader.Close()
			if readErr != nil {
				return nil, fmt.Errorf("failed to read gzip body: %w", readErr)
			} else if closeErr != nil {
				return nil, fmt.Errorf("failed to close gzip reader: %w", closeErr)
			} else {
				decodedBody = unzipped
			}
		}
	}

	return decodedBody, nil
}

func Parse(raw []byte) (*Response, error) {
	return ParseWithRedirectInfo(raw, 0)
}

func ParseWithRedirectInfo(raw []byte, redirectCount int) (*Response, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	if redirectCount < 0 {
		redirectCount = 0
	}

	statusLine, headerFields, rawBody, err := splitHeadersAndBody(raw)
	if err != nil {
		return nil, err
	}

	statusCode, err := parseStatusCode(statusLine)
	if err != nil {
		return nil, err
	}

	headersMap := parseHeaders(headerFields)
	decodedBody, err := decodeBodyByHeaders(headersMap, rawBody)
	if err != nil {
		return nil, err
	}

	contentType := ""
	if values := headersMap["Content-Type"]; len(values) > 0 {
		contentType = values[0]
	}

	return &Response{
		StatusLine:    statusLine,
		HeaderFields:  append([]string(nil), headerFields...),
		Body:          decodedBody,
		ResponseIsOK:  statusCode >= 200 && statusCode <= 299,
		IsRedirected:  redirectCount > 0,
		RedirectCount: redirectCount,
		ContentType:   contentType,
	}, nil
}
