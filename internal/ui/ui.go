package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go2web/internal/parser"
	"go2web/internal/search"
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
	Print("%s", formatParsedResponse(resp))
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
