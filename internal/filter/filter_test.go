package filter

import (
	"testing"
	"time"

	"github.com/atopos31/agent-rss/pkg/model"
)

func TestFilter_TimeRange(t *testing.T) {
	now := time.Now()
	items := []model.Item{
		{Title: "Old", Time: now.Add(-48 * time.Hour)},
		{Title: "Recent", Time: now.Add(-1 * time.Hour)},
		{Title: "Future", Time: now.Add(24 * time.Hour)},
	}

	since := now.Add(-24 * time.Hour)
	until := now

	filtered := Filter(items, Options{Since: &since, Until: &until})

	if len(filtered) != 1 {
		t.Fatalf("expected 1 item, got %d", len(filtered))
	}
	if filtered[0].Title != "Recent" {
		t.Fatalf("expected 'Recent', got %s", filtered[0].Title)
	}
}

func TestFilter_TitleKeyword(t *testing.T) {
	items := []model.Item{
		{Title: "Go Programming"},
		{Title: "Python Basics"},
		{Title: "Advanced Go Techniques"},
	}

	filtered := Filter(items, Options{Titles: []string{"go"}})

	if len(filtered) != 2 {
		t.Fatalf("expected 2 items, got %d", len(filtered))
	}
}

func TestFilter_ContentKeyword(t *testing.T) {
	items := []model.Item{
		{Title: "A", Content: "Learn about AI"},
		{Title: "B", Content: "Machine learning basics"},
		{Title: "C", Content: "Database design"},
	}

	filtered := Filter(items, Options{Contents: []string{"learning", "ai"}})

	if len(filtered) != 2 {
		t.Fatalf("expected 2 items, got %d", len(filtered))
	}
}

func TestFilter_Combined(t *testing.T) {
	now := time.Now()
	items := []model.Item{
		{Title: "Go News", Content: "Latest updates", Time: now.Add(-1 * time.Hour)},
		{Title: "Go News", Content: "Old updates", Time: now.Add(-48 * time.Hour)},
		{Title: "Python News", Content: "Latest updates", Time: now.Add(-1 * time.Hour)},
	}

	since := now.Add(-24 * time.Hour)
	filtered := Filter(items, Options{
		Since:  &since,
		Titles: []string{"go"},
	})

	if len(filtered) != 1 {
		t.Fatalf("expected 1 item, got %d", len(filtered))
	}
	if filtered[0].Title != "Go News" {
		t.Fatalf("expected 'Go News', got %s", filtered[0].Title)
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"2026-03-12T08:30:00Z", false},
		{"2026-03-12", false},
		{"invalid", true},
		{"2026/03/12", true},
	}

	for _, tt := range tests {
		_, err := ParseTime(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseTime(%q): wantErr=%v, got err=%v", tt.input, tt.wantErr, err)
		}
	}
}

func TestParseTime_Relative(t *testing.T) {
	tests := []struct {
		input    string
		minAgo   time.Duration
		maxAgo   time.Duration
		wantErr  bool
	}{
		{"1h", 59 * time.Minute, 61 * time.Minute, false},
		{"2h", 119 * time.Minute, 121 * time.Minute, false},
		{"1d", 23 * time.Hour, 25 * time.Hour, false},
		{"30m", 29 * time.Minute, 31 * time.Minute, false},
		{"0m", 0, 1 * time.Minute, false},
		{"abc", 0, 0, true},
		{"1x", 0, 0, true},
		{"h1", 0, 0, true},
	}

	for _, tt := range tests {
		parsed, err := ParseTime(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseTime(%q): expected error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseTime(%q): unexpected error: %v", tt.input, err)
			continue
		}

		ago := time.Since(parsed)
		if ago < tt.minAgo || ago > tt.maxAgo {
			t.Errorf("ParseTime(%q): got %v ago, want between %v and %v", tt.input, ago, tt.minAgo, tt.maxAgo)
		}
	}
}
