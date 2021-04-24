package color

import (
	"fmt"
)

// Escape characters
const (
	escape = "\x1b"
	reset  = "\x1b[0m"
)

// ANSI Color codes
const (
	black  = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

func colSprintf(col int, str string) string {
	return fmt.Sprintf("%s[%dm%s%s", escape, col, str, reset) 
}

func Red(str string) string {
	return colSprintf(red, str)
}

func Green(str string) string {
	return colSprintf(green, str)
}

func Yellow(str string) string {
	return colSprintf(yellow, str)
}

func Blue(str string) string {
	return colSprintf(blue, str)
}

func Magenta(str string) string {
	return colSprintf(magenta, str)
}

func Cyan(str string) string {
	return colSprintf(cyan, str)
}

func White(str string) string {
	return colSprintf(white, str)
}
