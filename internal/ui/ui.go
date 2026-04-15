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
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F8FAFC")).
			Background(lipgloss.Color("#1E293B")).
			Padding(0, 1)

	okBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#052E16")).
			Background(lipgloss.Color("#86EFAC")).
			Padding(0, 1)

	errBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#450A0A")).
			Background(lipgloss.Color("#FCA5A5")).
			Padding(0, 1)

	metaLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#475569"))

	metaValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0F172A"))

	headersBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#94A3B8")).
			Padding(0, 1)

	bodyPlaceholderStyle = lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color("#64748B")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#CBD5E1")).
				Padding(0, 1)
)

func Print(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(msg)
}

func Error(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func formatParsedResponse(resp *parser.Response) string {
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
	if resp == nil {
		lipgloss.Print(errBadgeStyle.Render("response: <nil>") + "\n")
		return
	}

	statusBadge := errBadgeStyle.Render("ERROR")
	if resp.ResponseIsOK {
		statusBadge = okBadgeStyle.Render("OK")
	}

	redirectValue := "no"
	if resp.IsRedirected {
		redirectValue = fmt.Sprintf("yes (%d)", resp.RedirectCount)
	}

	contentTypeStyle := metaValueStyle
	if strings.Contains(strings.ToLower(resp.ContentType), "text/html") {
		contentTypeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0369A1"))
	}

	headersText := "(none)"
	if len(resp.HeaderFields) > 0 {
		headersText = strings.Join(resp.HeaderFields, "\n")
	}

	statusLine := fmt.Sprintf(
		"%s %s\n%s",
		titleStyle.Render("URL Response"),
		statusBadge,
		resp.StatusLine,
	)

	metaLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		metaLabelStyle.Render("content type: ")+contentTypeStyle.Render(resp.ContentType),
		"   ",
		metaLabelStyle.Render("redirected: ")+metaValueStyle.Render(redirectValue),
	)

	headersBlock := lipgloss.JoinVertical(
		lipgloss.Left,
		metaLabelStyle.Render("headers"),
		headersBoxStyle.Render(headersText),
	)

	bodyBlock := lipgloss.JoinVertical(
		lipgloss.Left,
		metaLabelStyle.Render("body"),
		bodyPlaceholderStyle.Render("[body preview placeholder for upcoming renderer]"),
	)

	out := lipgloss.JoinVertical(
		lipgloss.Left,
		statusLine,
		metaLine,
		headersBlock,
		bodyBlock,
	)

	lipgloss.Print(out + "\n")
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
