package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/atopos31/agent-rss/pkg/model"
)

func TestWriter_JSON(t *testing.T) {
	items := []model.Item{
		{Name: "test", Title: "Title 1", Time: time.Date(2026, 3, 12, 8, 0, 0, 0, time.UTC)},
		{Name: "test", Title: "Title 2", Time: time.Date(2026, 3, 12, 9, 0, 0, 0, time.UTC)},
	}

	var buf bytes.Buffer
	w := New(&buf, FormatJSON)
	if err := w.Write(items); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	var result []model.Item
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
}

func TestWriter_NDJSON(t *testing.T) {
	items := []model.Item{
		{Name: "test", Title: "Title 1"},
		{Name: "test", Title: "Title 2"},
	}

	var buf bytes.Buffer
	w := New(&buf, FormatNDJSON)
	if err := w.Write(items); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	for i, line := range lines {
		var item model.Item
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			t.Fatalf("line %d: JSON unmarshal failed: %v", i, err)
		}
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input string
		want  Format
	}{
		{"json", FormatJSON},
		{"ndjson", FormatNDJSON},
		{"invalid", FormatNDJSON},
		{"", FormatNDJSON},
	}

	for _, tt := range tests {
		got := ParseFormat(tt.input)
		if got != tt.want {
			t.Errorf("ParseFormat(%q): want %v, got %v", tt.input, tt.want, got)
		}
	}
}
