package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Response struct {
	StatusLine   string
	HeaderFields []string
	Body         []byte
	ResponseIsOK bool
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

func Parse(raw []byte) (*Response, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	headerEnd := bytes.Index(raw, []byte("\r\n\r\n"))
	lineSeparator := "\r\n"
	bodyStart := 4

	if headerEnd == -1 {
		headerEnd = bytes.Index(raw, []byte("\n\n"))
		lineSeparator = "\n"
		bodyStart = 2
	}

	if headerEnd == -1 {
		return nil, fmt.Errorf("invalid http response: missing header separator")
	}

	lines := strings.Split(string(raw[:headerEnd]), lineSeparator)
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return nil, fmt.Errorf("invalid http response: missing status line")
	}

	statusCode, err := parseStatusCode(lines[0])
	if err != nil {
		return nil, err
	}

	headers := make([]string, 0, len(lines)-1)
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		headers = append(headers, line)
	}

	return &Response{
		StatusLine:   lines[0],
		HeaderFields: headers,
		Body:         raw[headerEnd+bodyStart:],
		ResponseIsOK: statusCode >= 200 && statusCode <= 299,
	}, nil
}
