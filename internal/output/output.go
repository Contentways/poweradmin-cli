// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output

// Format represents the output format.
type Format string

const (
	FormatTable Format = "table"
	FormatFull  Format = "full"
	FormatJSON  Format = "json"
)

// ParseFormat parses a format string.
func ParseFormat(s string) Format {
	switch s {
	case "full":
		return FormatFull
	case "json":
		return FormatJSON
	default:
		return FormatTable
	}
}
