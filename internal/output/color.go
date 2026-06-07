// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output

import (
	"os"

	"golang.org/x/term"
)

// ANSI escape codes for terminal colors and formatting.
const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	gray   = "\033[90m"
)

// IsTTY returns true if stdout is a terminal.
// Used to suppress ANSI codes when output is piped or redirected.
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Bold wraps s in bold ANSI codes if stdout is a TTY.
func Bold(s string) string {
	if !IsTTY() {
		return s
	}
	return bold + s + reset
}

// Cyan wraps s in cyan ANSI codes if stdout is a TTY.
func Cyan(s string) string {
	if !IsTTY() {
		return s
	}
	return cyan + s + reset
}

// Green wraps s in green ANSI codes if stdout is a TTY.
func Green(s string) string {
	if !IsTTY() {
		return s
	}
	return green + s + reset
}

// Red wraps s in red ANSI codes if stdout is a TTY.
func Red(s string) string {
	if !IsTTY() {
		return s
	}
	return red + s + reset
}

// Yellow wraps s in yellow ANSI codes if stdout is a TTY.
func Yellow(s string) string {
	if !IsTTY() {
		return s
	}
	return yellow + s + reset
}

// Gray wraps s in gray ANSI codes if stdout is a TTY.
func Gray(s string) string {
	if !IsTTY() {
		return s
	}
	return gray + s + reset
}
