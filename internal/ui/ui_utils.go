package ui

import (
	"bytes"
	"encoding/json"
	"strings"

	"go2web/internal/parser"

	"golang.org/x/net/html"
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

type textBuilder struct {
	builder        strings.Builder
	lastWasSpace   bool
	lastWasNewline bool
}

func (t *textBuilder) writeNewline() {
	if t.builder.Len() == 0 || t.lastWasNewline {
		return
	}

	t.builder.WriteByte('\n')
	t.lastWasNewline = true
	t.lastWasSpace = false
}

func (t *textBuilder) writeWord(word string) {
	if word == "" {
		return
	}

	if t.builder.Len() > 0 && !t.lastWasSpace && !t.lastWasNewline {
		t.builder.WriteByte(' ')
	}

	t.builder.WriteString(word)
	t.lastWasSpace = false
	t.lastWasNewline = false
}

func (t *textBuilder) writeText(text string) {
	for _, field := range strings.Fields(text) {
		t.writeWord(field)
	}
}

func (t *textBuilder) writeBullet() {
	if t.builder.Len() > 0 {
		t.writeNewline()
	}

	t.builder.WriteString("- ")
	t.lastWasSpace = true
	t.lastWasNewline = false
}

func (t *textBuilder) String() string {
	return t.builder.String()
}

func normalizeWhitespace(text string) string {
	lines := strings.Split(text, "\n")
	cleaned := make([]string, 0, len(lines))
	emptyLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			emptyLines++
			if emptyLines > 1 {
				continue
			}
			cleaned = append(cleaned, "")
			continue
		}

		emptyLines = 0
		cleaned = append(cleaned, trimmed)
	}

	return strings.Join(cleaned, "\n")
}

func htmlToText(body []byte) (string, error) {
	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	builder := &textBuilder{}
	walkHTMLNodes(root, builder)
	return normalizeWhitespace(builder.String()), nil
}

func walkHTMLNodes(node *html.Node, builder *textBuilder) {
	if node == nil {
		return
	}

	if node.Type == html.ElementNode {
		tag := strings.ToLower(node.Data)
		switch tag {
		case "script", "style", "noscript", "svg", "canvas", "head", "meta", "link":
			return
		case "br", "hr":
			builder.writeNewline()
			return
		case "li":
			builder.writeBullet()
		case "p", "div", "section", "article", "header", "footer", "nav", "main", "aside", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "table", "tr", "td", "th", "pre", "blockquote":
			builder.writeNewline()
		}
	}

	if node.Type == html.TextNode {
		builder.writeText(node.Data)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		walkHTMLNodes(child, builder)
	}

	if node.Type == html.ElementNode {
		tag := strings.ToLower(node.Data)
		switch tag {
		case "p", "div", "section", "article", "header", "footer", "nav", "main", "aside", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "table", "tr", "td", "th", "pre", "blockquote":
			builder.writeNewline()
		}
	}
}

func prettyPrintJSON(body []byte) string {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return ""
	}

	var out bytes.Buffer
	if err := json.Indent(&out, trimmed, "", "  "); err != nil {
		return string(body)
	}

	return out.String()
}

func renderBody(resp *parser.Response) string {
	if resp == nil {
		return "(no response body)"
	}

	if len(resp.Body) == 0 {
		return "(empty body)"
	}

	contentType := strings.ToLower(resp.ContentType)
	var bodyText string

	switch {
	case strings.Contains(contentType, "json"):
		bodyText = prettyPrintJSON(resp.Body)
	case strings.Contains(contentType, "html"):
		text, err := htmlToText(resp.Body)
		if err != nil {
			bodyText = string(resp.Body)
		} else {
			bodyText = text
		}
	default:
		bodyText = string(resp.Body)
	}

	bodyText = strings.TrimSpace(bodyText)
	if bodyText == "" {
		return "(empty body)"
	}

	return ensureMaxLineLength(bodyText, maxResponseLineLength)
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

	normalized.Body = resp.Body

	return &normalized
}
