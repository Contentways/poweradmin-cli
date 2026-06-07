// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package output provides shared output formatting utilities for the CLI.
// It supports three output formats: table (default), full and json.
// Color output is automatically suppressed when stdout is not a terminal
// (e.g. when piping output to another command or redirecting to a file).
package output

import (
	"os"

	"golang.org/x/term"
)

// ANSI escape codes for terminal colors and formatting.
// These are only applied when stdout is a TTY — see IsTTY().
const (
	reset  = "\033[0m"  // Resets all attributes to default.
	bold   = "\033[1m"  // Bold text — used for table headers.
	cyan   = "\033[36m" // Cyan — used for record and zone types.
	green  = "\033[32m" // Green — used for active status and non-zero counts.
	red    = "\033[31m" // Red — used for inactive status.
	yellow = "\033[33m" // Yellow — reserved for warnings.
	gray   = "\033[90m" // Gray — reserved for secondary information.
)

// IsTTY returns true if stdout is connected to a terminal.
// When false, all color functions return their input unchanged,
// ensuring clean output when piping or redirecting.
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Bold wraps s in bold ANSI formatting if stdout is a TTY.
// Used for table headers to visually distinguish them from data rows.
func Bold(s string) string {
	if !IsTTY() {
		return s
	}
	return bold + s + reset
}

// Cyan wraps s in cyan ANSI color if stdout is a TTY.
// Used to highlight DNS record types (A, AAAA, MX, TXT, etc.)
// and zone types (NATIVE, MASTER, SLAVE).
func Cyan(s string) string {
	if !IsTTY() {
		return s
	}
	return cyan + s + reset
}

// Green wraps s in green ANSI color if stdout is a TTY.
// Used to highlight positive states such as active users
// and non-zero member or zone counts in groups.
func Green(s string) string {
	if !IsTTY() {
		return s
	}
	return green + s + reset
}

// Red wraps s in red ANSI color if stdout is a TTY.
// Used to highlight negative states such as inactive users.
func Red(s string) string {
	if !IsTTY() {
		return s
	}
	return red + s + reset
}

// Yellow wraps s in yellow ANSI color if stdout is a TTY.
// Reserved for warning states and advisory output.
func Yellow(s string) string {
	if !IsTTY() {
		return s
	}
	return yellow + s + reset
}

// Gray wraps s in gray ANSI color if stdout is a TTY.
// Reserved for secondary or de-emphasized information.
func Gray(s string) string {
	if !IsTTY() {
		return s
	}
	return gray + s + reset
}

// CyanCode returns the raw ANSI escape code for cyan.
// Used with Cell() in AddColoredRow to colorize table cells
// after tabwriter alignment, avoiding column width distortion.
func CyanCode() string { return cyan }

// GreenCode returns the raw ANSI escape code for green.
// Used with Cell() in AddColoredRow to colorize table cells
// after tabwriter alignment, avoiding column width distortion.
func GreenCode() string { return green }

// RedCode returns the raw ANSI escape code for red.
// Used with Cell() in AddColoredRow to colorize table cells
// after tabwriter alignment, avoiding column width distortion.
func RedCode() string { return red }
