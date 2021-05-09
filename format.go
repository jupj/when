package main

import (
	"fmt"
	"strings"
)

const (
	// VT100 escape sequences for formatting
	vtReset   = "0"
	vtBold    = "1"
	vtReverse = "7"
	vtGreen   = "32"
)

// escape returns a VT100 escape sequence with the given codes
func escape(code ...string) string {
	return fmt.Sprintf("\x1b[%sm", strings.Join(code, ";"))
}

// colFmt formats hour strings with style/color
func colFmt(s string, h int, currentHour bool) string {
	var codes []string

	// Highlight current hour
	if currentHour {
		codes = append(codes, vtReverse)
	}

	// Bold the first hour of the day
	if h == 0 {
		codes = append(codes, vtBold)
	}

	// Green color for office hours
	if 8 <= h && h <= 17 {
		codes = append(codes, vtGreen)
	}

	if len(codes) == 0 {
		return s
	}

	return escape(codes...) + s + escape(vtReset)
}
