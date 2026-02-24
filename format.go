package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	// VT100 escape sequences for formatting
	vtReset   = "0"
	vtBold    = "1"
	vtReverse = "7"
	vtRed     = "31"
	vtGreen   = "32"
	vtBlue    = "34"
)

// escape returns a VT100 escape sequence with the given codes
func escape(code ...string) string {
	return fmt.Sprintf("\x1b[%sm", strings.Join(code, ";"))
}

// dstChanged checks for changes in daylight saving time.
// Returns true if [t +/- 1h] doesn't differ 1h on the "wall clock".
func dstChanged(t time.Time) bool {
	h := t.Hour()

	prev := t.Add(-time.Hour).Hour()
	expectedPrev := (h + 23) % 24
	if prev != expectedPrev {
		return true
	}

	next := t.Add(time.Hour).Hour()
	expectedNext := (h + 1) % 24

	return next != expectedNext
}

// colFmt formats hour strings with style/color
func colFmt(s string, t time.Time, currentHour bool) string {
	var codes []string
	h := t.Hour()

	// Highlight current hour
	if currentHour {
		codes = append(codes, vtReverse)
	}

	// Bold the first hour of the day
	if h == 0 {
		codes = append(codes, vtBold)
	}

	// Coloring - use separate colors for each 8-hour block
	switch {
	case dstChanged(t):
		// Highlight change in daylight saving
		codes = append(codes, vtRed)
	case 8 <= h && h < 16:
		// Green color for day hours
		codes = append(codes, vtGreen)
	case h >= 16:
		// Blue color for evening hours
		codes = append(codes, vtBlue)
	}

	if len(codes) == 0 {
		return s
	}

	return escape(codes...) + s + escape(vtReset)
}
