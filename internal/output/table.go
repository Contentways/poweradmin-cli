// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output

import (
	"io"
	"text/tabwriter"
)

// Table wraps a tabwriter.Writer to render aligned, padded text tables.
// Columns are separated by tab characters and automatically aligned
// based on the widest value in each column.
type Table struct {
	w        *tabwriter.Writer
	noHeader bool
}

// New creates a new Table that writes to w.
// The tabwriter is configured with 3 spaces of padding between columns.
func New(w io.Writer) *Table {
	return &Table{w: tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)}
}

// SetNoHeader controls whether the header row is rendered.
// When set to true, calls to AddHeader are silently ignored.
// Must be called before AddHeader to take effect.
func (t *Table) SetNoHeader(v bool) {
	t.noHeader = v
}

// AddHeader writes a header row to the table.
// Column values are separated by tab characters for alignment.
// If SetNoHeader(true) was called, this method is a no-op.
func (t *Table) AddHeader(columns ...string) {
	if t.noHeader {
		return
	}
	for i, col := range columns {
		if i > 0 {
			t.w.Write([]byte("\t"))
		}
		t.w.Write([]byte(col))
	}
	t.w.Write([]byte("\n"))
}

// AddRow writes a data row to the table.
// Column values are separated by tab characters for alignment.
func (t *Table) AddRow(columns ...string) {
	for i, col := range columns {
		if i > 0 {
			t.w.Write([]byte("\t"))
		}
		t.w.Write([]byte(col))
	}
	t.w.Write([]byte("\n"))
}

// Flush flushes all buffered data to the underlying writer.
// Must be called after all rows have been added to render the table.
func (t *Table) Flush() {
	t.w.Flush()
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
