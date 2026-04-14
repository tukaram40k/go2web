package ui

import (
	"fmt"
	"os"

	"go2web/internal/parser"
)

func Print(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(msg)
}

func Error(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func PrintParsedResponse(resp *parser.Response) {
	if resp == nil {
		Print("response: <nil>\n")
		return
	}

	Print("status line:\n%s\n\n", resp.StatusLine)
	Print("headers:\n")
	if len(resp.HeaderFields) == 0 {
		Print("(none)\n")
	} else {
		for _, header := range resp.HeaderFields {
			Print("%s\n", header)
		}
	}

	Print("\nbody:\n%s\n", string(resp.Body))
}

func Log(data []byte) (string, error) {
	logDir := ".go2web"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", err
	}

	f, err := os.CreateTemp(logDir, "go2web-response-*.txt")
	if err != nil {
		return "", err
	}

	if _, err := f.Write(data); err != nil {
		f.Close()
		return "", err
	}

	if err := f.Close(); err != nil {
		return "", err
	}

	return f.Name(), nil
}
