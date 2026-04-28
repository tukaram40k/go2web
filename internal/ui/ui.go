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

func Print(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(msg)
}

func Error(format string, a ...any) error {
	return fmt.Errorf(format, a...)
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

	bodyText := renderBody(resp)
	msg += fmt.Sprintf("\nbody:\n%s\n", bodyText)

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

	bodyText := renderBody(resp)
	bodyBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		bodyTextStyle.Render(bodyText),
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
		lipgloss.WithWhitespaceChars("."),
		lipgloss.WithWhitespaceStyle(backgroundStyle),
	)

	lipgloss.Print(rendered + "\n")
}

func PrintSearchResults(term string, resp *parser.Response, results []search.Result) {
	if resp != nil {
		resp = normalizedResponseLines(resp, maxResponseLineLength)
	}

	statusBadge := errBadgeStyle.Render("ERROR")
	if resp != nil && resp.ResponseIsOK {
		statusBadge = okBadgeStyle.Render("OK")
	}

	searchStatusBlockStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())
	if resp != nil && resp.ResponseIsOK {
		searchStatusBlockStyle = searchStatusBlockStyle.
			BorderForeground(lipgloss.Color(colorBorderOK))
	} else {
		searchStatusBlockStyle = searchStatusBlockStyle.
			BorderForeground(lipgloss.Color(colorBorderError))
	}

	searchTitle := "Search Results"
	if strings.TrimSpace(term) != "" {
		searchTitle = fmt.Sprintf("Search Results: %s", term)
	}

	searchStatusBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, searchTitle, " ", statusBadge),
	)
	searchStatusBlock = panelStyle.Render(searchStatusBlock)

	statusValue := "unknown"
	contentTypeValue := "unknown"
	redirectCountValue := "0"

	if resp != nil {
		statusValue = resp.StatusLine
		contentTypeValue = resp.ContentType
		redirectCountValue = strconv.Itoa(resp.RedirectCount)
		if resp.IsRedirected {
			redirectCountValue = redirectCountValue + " (redirected)"
		}
	}

	resultCountValue := strconv.Itoa(len(results))

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
			[]string{"content type", metaValueStyle.Render(contentTypeValue)},
			[]string{"redirect count", redirectCountValue},
			[]string{"results", resultCountValue},
		)

	resultsText := formatSearchResults(results)
	resultsBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		resultsText,
	)
	resultsBlock = panelStyle.Render(resultsBlock)

	metaBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		metaTable.String(),
	)
	metaBlock = panelStyle.Render(metaBlock)

	sharedContentWidth := lipgloss.Width(resultsBlock)
	if w := lipgloss.Width(searchStatusBlock); w > sharedContentWidth {
		sharedContentWidth = w
	}

	stretchStyle := lipgloss.NewStyle().Width(sharedContentWidth).Align(lipgloss.Center)
	searchStatusBlock = stretchStyle.Render(searchStatusBlock)
	resultsBlock = stretchStyle.Render(resultsBlock)

	searchStatusBlock = searchStatusBlockStyle.Render(searchStatusBlock)
	resultsBlock = searchResultsBlockStyle.Render(resultsBlock)

	out := strings.Join([]string{
		searchStatusBlock,
		metaBlock,
		resultsBlock,
	}, "\n\n")

	canvasWidth := lipgloss.Width(out) + 12
	if canvasWidth < 96 {
		canvasWidth = 96
	}
	canvasHeight := lipgloss.Height(out) + 4
	if canvasHeight < 24 {
		canvasHeight = 24
	}

	backgroundStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorTextBackground))

	rendered := lipgloss.Place(
		canvasWidth,
		canvasHeight,
		lipgloss.Center,
		lipgloss.Center,
		out,
		lipgloss.WithWhitespaceChars("."),
		lipgloss.WithWhitespaceStyle(backgroundStyle),
	)

	lipgloss.Print(rendered + "\n")
}

func formatSearchResults(results []search.Result) string {
	if len(results) == 0 {
		return searchResultTitleStyle.Render("no results found")
	}

	var builder strings.Builder

	for i, r := range results {
		title := strings.TrimSpace(r.Title)
		urlValue := strings.TrimSpace(r.URL)

		if title == "" {
			title = "(untitled)"
		}

		title = ensureMaxLineLength(title, maxResponseLineLength)
		urlValue = ensureMaxLineLength(urlValue, maxResponseLineLength)

		if i > 0 {
			builder.WriteString("\n")
		}

		builder.WriteString(fmt.Sprintf("%d. ", i+1))
		builder.WriteString(searchResultTitleStyle.Render(title))

		if urlValue != "" {
			builder.WriteString("\n   ")
			builder.WriteString(searchResultURLStyle.Render(urlValue))
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func formatSearchLog(term string, resp *parser.Response, results []search.Result) string {
	var builder strings.Builder
	if resp != nil {
		builder.WriteString(formatParsedResponse(resp))
		builder.WriteString("\n")
	}

	term = strings.TrimSpace(term)
	if term != "" {
		builder.WriteString(fmt.Sprintf("search term: %s\n", term))
	}

	builder.WriteString("\nresults:\n")
	if len(results) == 0 {
		builder.WriteString("no results found\n")
		return builder.String()
	}

	for i, r := range results {
		title := strings.TrimSpace(r.Title)
		urlValue := strings.TrimSpace(r.URL)

		if title == "" {
			title = "(untitled)"
		}

		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, title))
		if urlValue != "" {
			builder.WriteString(fmt.Sprintf("   %s\n", urlValue))
		}
	}

	return builder.String()
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

func LogSearchResults(term string, resp *parser.Response, results []search.Result) (string, error) {
	logDir := ".go2web"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", err
	}

	fileName := "search-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".txt"
	filePath := filepath.Join(logDir, fileName)

	logText := formatSearchLog(term, resp, results)
	if err := os.WriteFile(filePath, []byte(logText), 0o644); err != nil {
		return "", err
	}

	return filePath, nil
}
