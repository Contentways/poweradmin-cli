// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package output

import (
	"io"
	"text/tabwriter"
)

// Table is a simple tabwriter-based table output.
type Table struct {
	w *tabwriter.Writer
}

// New creates a new Table writing to w.
func New(w io.Writer) *Table {
	return &Table{
		w: tabwriter.NewWriter(w, 0, 0, 3, ' ', 0),
	}
}

// AddHeader writes a header row.
func (t *Table) AddHeader(columns ...string) {
	for i, col := range columns {
		if i > 0 {
			t.w.Write([]byte("\t"))
		}
		t.w.Write([]byte(col))
	}
	t.w.Write([]byte("\n"))
}

// AddRow writes a data row.
func (t *Table) AddRow(columns ...string) {
	for i, col := range columns {
		if i > 0 {
			t.w.Write([]byte("\t"))
		}
		t.w.Write([]byte(col))
	}
	t.w.Write([]byte("\n"))
}

// Truncate shortens s to max characters and adds "..." if truncated.
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Flush flushes the table to the writer.
func (t *Table) Flush() {
	t.w.Flush()
}
