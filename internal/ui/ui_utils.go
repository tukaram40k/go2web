package ui

import (
	"strings"

	"go2web/internal/parser"
)

func ensureMaxLineLength(text string, maxChars int) string {
	if maxChars <= 0 || text == "" {
		return text
	}

	lines := strings.Split(text, "\n")
	wrapped := make([]string, 0, len(lines))

	for _, line := range lines {
		runes := []rune(line)
		if len(runes) <= maxChars {
			wrapped = append(wrapped, line)
			continue
		}

		for start := 0; start < len(runes); start += maxChars {
			end := start + maxChars
			if end > len(runes) {
				end = len(runes)
			}
			wrapped = append(wrapped, string(runes[start:end]))
		}
	}

	return strings.Join(wrapped, "\n")
}

func normalizedResponseLines(resp *parser.Response, maxChars int) *parser.Response {
	if resp == nil {
		return nil
	}

	normalized := *resp
	normalized.StatusLine = ensureMaxLineLength(resp.StatusLine, maxChars)
	normalized.ContentType = ensureMaxLineLength(resp.ContentType, maxChars)

	normalized.HeaderFields = make([]string, len(resp.HeaderFields))
	for i, header := range resp.HeaderFields {
		normalized.HeaderFields[i] = ensureMaxLineLength(header, maxChars)
	}

	normalized.Body = []byte(ensureMaxLineLength(string(resp.Body), maxChars))

	return &normalized
}