package parser

import (
	"bytes"
	"fmt"
	"strings"
)

type Response struct {
	StatusLine   string
	HeaderFields []string
	Body         []byte
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
	}, nil
}
