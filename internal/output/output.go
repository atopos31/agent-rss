// Package output provides JSON and NDJSON serialization for items.
package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/atopos31/agent-rss/pkg/model"
)

// Format represents the output format.
type Format string

const (
	FormatJSON   Format = "json"
	FormatNDJSON Format = "ndjson"
)

// Writer outputs items in the specified format.
type Writer struct {
	w      io.Writer
	format Format
}

// New creates a new Writer.
func New(w io.Writer, format Format) *Writer {
	return &Writer{w: w, format: format}
}

// Write outputs all items in the configured format.
func (w *Writer) Write(items []model.Item) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(items)
	case FormatNDJSON:
		return w.writeNDJSON(items)
	default:
		return w.writeNDJSON(items)
	}
}

func (w *Writer) writeJSON(items []model.Item) error {
	enc := json.NewEncoder(w.w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(items); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func (w *Writer) writeNDJSON(items []model.Item) error {
	enc := json.NewEncoder(w.w)
	for _, item := range items {
		if err := enc.Encode(item); err != nil {
			return fmt.Errorf("encode ndjson: %w", err)
		}
	}
	return nil
}

// ParseFormat parses a format string.
func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "ndjson":
		return FormatNDJSON
	default:
		return FormatNDJSON
	}
}
