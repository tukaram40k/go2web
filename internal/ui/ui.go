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
