package ui

import (
	"errors"
	"fmt"
)

func Print(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(msg)
}

func Error(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return errors.New(msg)
}
