// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// coloredCell holds a plain text value and an optional ANSI color code.
// The plain value is used for tabwriter alignment; the color is applied
// after the table has been rendered so escape codes do not affect column widths.
type coloredCell struct {
	plain string
	color string
}

// Cell creates a coloredCell with the given plain text value and ANSI color code.
// Use the exported color code functions (CyanCode, GreenCode, etc.) for the color.
// Pass an empty string for color to render the cell without color.
func Cell(plain, color string) coloredCell {
	return coloredCell{plain: plain, color: color}
}

// PlainCell creates a coloredCell with no color — the value is rendered as-is.
func PlainCell(plain string) coloredCell {
	return coloredCell{plain: plain}
}

// Table wraps a tabwriter.Writer to render aligned, padded text tables.
// Colors are stored separately from plain text values and applied after
// tabwriter alignment so ANSI escape codes do not affect column widths.
type Table struct {
	// buf holds the plain text output from tabwriter before colorization.
	buf bytes.Buffer
	// w is the tabwriter that aligns columns using plain text values.
	w *tabwriter.Writer
	// finalWriter is the destination writer for the colorized output.
	finalWriter io.Writer
	// noHeader suppresses the header row when true.
	noHeader bool
	// header holds the plain text header columns.
	header []string
	// rows holds the plain text data rows for tabwriter alignment.
	rows [][]string
	// colors holds per-cell ANSI color codes, parallel to rows.
	// An empty string means no color for that cell.
	colors [][]string
}

// New creates a new Table that writes its final output to w.
// The tabwriter is configured with 3 spaces of padding between columns.
func New(w io.Writer) *Table {
	t := &Table{finalWriter: w}
	t.w = tabwriter.NewWriter(&t.buf, 0, 0, 3, ' ', 0)
	return t
}

// SetNoHeader controls whether the header row is rendered.
// When set to true, calls to AddHeader are silently ignored.
// Must be called before AddHeader to take effect.
func (t *Table) SetNoHeader(v bool) {
	t.noHeader = v
}

// AddHeader stores the header column labels.
// The header is rendered in bold when the output is a TTY.
// If SetNoHeader(true) was called, this method is a no-op.
func (t *Table) AddHeader(columns ...string) {
	if t.noHeader {
		return
	}
	t.header = columns
}

// AddRow stores a plain text data row with no per-cell colors.
// Use AddColoredRow to specify colors for individual cells.
func (t *Table) AddRow(columns ...string) {
	t.rows = append(t.rows, columns)
	colors := make([]string, len(columns))
	t.colors = append(t.colors, colors)
}

// AddColoredRow stores a data row where each cell can carry an optional
// ANSI color code. Colors are applied after tabwriter alignment so they
// do not affect column widths. Use Cell() and PlainCell() to build cells.
func (t *Table) AddColoredRow(cells ...coloredCell) {
	plain := make([]string, len(cells))
	colors := make([]string, len(cells))
	for i, c := range cells {
		plain[i] = c.plain
		colors[i] = c.color
	}
	t.rows = append(t.rows, plain)
	t.colors = append(t.colors, colors)
}

// Flush writes all buffered rows to the tabwriter for column alignment,
// then applies per-cell ANSI colors and writes the final output to the
// underlying writer. Must be called after all rows have been added.
func (t *Table) Flush() {
	// Feed plain text to tabwriter for alignment.
	if len(t.header) > 0 {
		fmt.Fprintln(t.w, strings.Join(t.header, "\t"))
	}
	for _, row := range t.rows {
		fmt.Fprintln(t.w, strings.Join(row, "\t"))
	}
	t.w.Flush()

	// Apply colors line by line on the aligned output.
	lines := strings.Split(t.buf.String(), "\n")
	lineIdx := 0

	// Header line — render in bold when output is a TTY.
	if len(t.header) > 0 && lineIdx < len(lines) {
		if IsTTY() {
			fmt.Fprintln(t.finalWriter, bold+lines[lineIdx]+reset)
		} else {
			fmt.Fprintln(t.finalWriter, lines[lineIdx])
		}
		lineIdx++
	}

	// Data rows — apply per-cell colors by replacing plain values in the
	// aligned line with their colorized equivalents.
	for rowIdx, row := range t.rows {
		if lineIdx >= len(lines) {
			break
		}
		line := lines[lineIdx]
		lineIdx++

		// Skip colorization when not a TTY or no colors defined for this row.
		if !IsTTY() || rowIdx >= len(t.colors) {
			fmt.Fprintln(t.finalWriter, line)
			continue
		}

		colored := line
		rowColors := t.colors[rowIdx]
		for colIdx, plain := range row {
			if colIdx < len(rowColors) && rowColors[colIdx] != "" && plain != "" {
				colored = strings.Replace(colored, plain, rowColors[colIdx]+plain+reset, 1)
			}
		}
		fmt.Fprintln(t.finalWriter, colored)
	}
}

// Truncate shortens s to at most max characters.
// If s exceeds max, it is cut at max-3 characters and "..." is appended.
// This keeps long values (e.g. TXT record content) from breaking table alignment.
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
