package ui

import (
	"fmt"
	"os"
)

func Print(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(msg)
}

func Error(format string, a ...any) error {
	return fmt.Errorf(format, a...)
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
