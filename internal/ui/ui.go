package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go2web/internal/parser"
	"go2web/internal/search"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var (
	okBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextOKBadge)).
			Background(lipgloss.Color("#86EFAC")).
			Padding(0, 1)

	errBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextErrorBadge)).
			Background(lipgloss.Color("#FCA5A5")).
			Padding(0, 1)

	metaLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextPrimary))

	metaValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	panelStyle = lipgloss.NewStyle().
			Margin(0, 0, 0, 0).
			Padding(0, 1)

	tableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	tableOddRowStyle = tableCellStyle.
				Foreground(lipgloss.Color(colorTextPrimary))

	tableEvenRowStyle = tableCellStyle.
				Foreground(lipgloss.Color(colorTextSecondary))

	headersBoxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary)).
			Padding(0, 1)

	headersBlockStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(colorBorderPrimary))

	bodyPlaceholderStyle = lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color(colorTextPrimary))

	bodyBlockStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorBorderPrimary))
)

const (
	maxResponseLineLength = 80

	colorTextPrimary    = "#b7d0f8"
	colorTextSecondary  = "#b7f8f4"
	colorTextOKBadge    = "#052E16"
	colorTextErrorBadge = "#450A0A"
	colorTextHTML       = "#0369A1"
	colorTextBackground = "#5b435a"

	colorBorderPrimary = "#8921d8"
	colorBorderOK      = "#15803D"
	colorBorderError   = "#B91C1C"
)

func Print(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(msg)
}

func Error(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

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

func formatParsedResponse(resp *parser.Response) string {
	resp = normalizedResponseLines(resp, maxResponseLineLength)

	if resp == nil {
		return "response: <nil>\n"
	}

	msg := ""
	msg += fmt.Sprintf("status line:\n%s\n", resp.StatusLine)
	msg += fmt.Sprintf("response ok: %t\n", resp.ResponseIsOK)
	msg += fmt.Sprintf("content type: %s\n", resp.ContentType)
	msg += fmt.Sprintf("redirected: %t\n", resp.IsRedirected)
	msg += fmt.Sprintf("redirect count: %d\n\n", resp.RedirectCount)

	msg += "headers:\n"
	if len(resp.HeaderFields) == 0 {
		msg += "(none)\n"
	} else {
		for _, header := range resp.HeaderFields {
			msg += fmt.Sprintf("%s\n", header)
		}
	}

	msg += fmt.Sprintf("\nbody:\n%s\n", string(resp.Body))

	return msg
}

func PrintParsedResponse(resp *parser.Response) {
	resp = normalizedResponseLines(resp, maxResponseLineLength)

	if resp == nil {
		lipgloss.Print(errBadgeStyle.Render("response: <nil>") + "\n")
		return
	}

	statusBadge := errBadgeStyle.Render("ERROR")
	if resp.ResponseIsOK {
		statusBadge = okBadgeStyle.Render("OK")
	}

	headersText := "(none)"
	if len(resp.HeaderFields) > 0 {
		headersText = strings.Join(resp.HeaderFields, "\n")
	}

	statusValue := resp.StatusLine

	urlResponseBlockStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())
	if resp.ResponseIsOK {
		urlResponseBlockStyle = urlResponseBlockStyle.
			BorderForeground(lipgloss.Color(colorBorderOK))
	} else {
		urlResponseBlockStyle = urlResponseBlockStyle.
			BorderForeground(lipgloss.Color(colorBorderError))
	}

	redirectCountValue := strconv.Itoa(resp.RedirectCount)
	if resp.IsRedirected {
		redirectCountValue = redirectCountValue + " (redirected)"
	}

	metaTable := table.New().
		Border(lipgloss.NormalBorder()).
		BorderHeader(false).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(colorBorderPrimary))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row%2 == 0:
				return tableEvenRowStyle
			default:
				return tableOddRowStyle
			}
		}).
		Rows(
			[]string{"status line", statusValue},
			[]string{"content type", metaValueStyle.Render(resp.ContentType)},
			[]string{"redirect count", redirectCountValue},
		)

	urlStatusBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, "URL Response", " ", statusBadge),
	)
	urlStatusBlock = panelStyle.Render(urlStatusBlock)

	tableBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		metaTable.String(),
	)
	tableBlock = panelStyle.Render(tableBlock)

	headersBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		headersBoxStyle.Render(headersText),
	)
	headersBlock = panelStyle.Render(headersBlock)

	bodyBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		bodyPlaceholderStyle.Render("[body preview placeholder for upcoming renderer]"),
	)
	bodyBlock = panelStyle.Render(bodyBlock)

	sharedContentWidth := lipgloss.Width(headersBlock)
	if w := lipgloss.Width(bodyBlock); w > sharedContentWidth {
		sharedContentWidth = w
	}

	stretchStyle := lipgloss.NewStyle().Width(sharedContentWidth).Align(lipgloss.Center)
	headersBlock = stretchStyle.Render(headersBlock)
	bodyBlock = stretchStyle.Render(bodyBlock)
	urlStatusBlock = stretchStyle.Render(urlStatusBlock)
	headersBlock = headersBlockStyle.Render(headersBlock)
	bodyBlock = bodyBlockStyle.Render(bodyBlock)
	urlStatusBlock = urlResponseBlockStyle.Render(urlStatusBlock)

	out := strings.Join([]string{
		urlStatusBlock,
		tableBlock,
		headersBlock,
		bodyBlock,
	}, "\n\n")

	canvasWidth := lipgloss.Width(out) + 12
	if canvasWidth < 96 {
		canvasWidth = 96
	}
	canvasHeight := lipgloss.Height(out) + 4
	if canvasHeight < 28 {
		canvasHeight = 28
	}

	backgroundStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextBackground))

	rendered := lipgloss.Place(
		canvasWidth,
		canvasHeight,
		lipgloss.Center,
		lipgloss.Center,
		out,
		lipgloss.WithWhitespaceChars("{}"),
		lipgloss.WithWhitespaceStyle(backgroundStyle),
	)

	lipgloss.Print(rendered + "\n")
}

func PrintSearchResults(term string, resp *parser.Response, results []search.Result) {
	Print("search term: %s\n", term)
	if resp != nil {
		Print("status line: %s\n", resp.StatusLine)
		Print("response ok: %t\n", resp.ResponseIsOK)
		Print("content type: %s\n", resp.ContentType)
		Print("redirected: %t\n", resp.IsRedirected)
		Print("redirect count: %d\n", resp.RedirectCount)
	}

	Print("\ntop results:\n")
	if len(results) == 0 {
		Print("no results found\n")
		return
	}

	for i, r := range results {
		Print("%d. %s\n", i+1, r.Title)
		Print("   %s\n", r.URL)
	}
}

func Log(resp *parser.Response) (string, error) {
	logDir := ".go2web"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", err
	}

	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ".txt"
	filePath := filepath.Join(logDir, fileName)

	if err := os.WriteFile(filePath, []byte(formatParsedResponse(resp)), 0o644); err != nil {
		return "", err
	}

	return filePath, nil
}
