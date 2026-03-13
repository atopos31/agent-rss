// Package filter provides time and keyword filtering for RSS items.
package filter

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/atopos31/agent-rss/pkg/model"
)

// Options defines filtering criteria.
type Options struct {
	Since    *time.Time
	Until    *time.Time
	Titles   []string
	Contents []string
}

// Filter returns items matching all specified criteria.
// An empty Options matches all items.
func Filter(items []model.Item, opts Options) []model.Item {
	result := make([]model.Item, 0, len(items))

	for _, item := range items {
		if matches(item, opts) {
			result = append(result, item)
		}
	}

	return result
}

func matches(item model.Item, opts Options) bool {
	if opts.Since != nil && item.Time.Before(*opts.Since) {
		return false
	}

	if opts.Until != nil && item.Time.After(*opts.Until) {
		return false
	}

	if len(opts.Titles) > 0 && !containsAny(item.Title, opts.Titles) {
		return false
	}

	if len(opts.Contents) > 0 && !containsAny(item.Content, opts.Contents) {
		return false
	}

	return true
}

func containsAny(text string, keywords []string) bool {
	lower := strings.ToLower(text)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// ParseTime attempts to parse a time string.
// Supported formats:
//   - RFC3339: 2026-03-12T08:30:00Z
//   - Date: 2026-03-12
//   - Relative: 1h, 2d, 30m (hours, days, minutes ago)
func ParseTime(s string) (time.Time, error) {
	// Try RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// Try YYYY-MM-DD
	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		return t, nil
	}

	// Try relative time (e.g., 1h, 2d, 30m)
	if t, ok := parseRelativeTime(s); ok {
		return t, nil
	}

	return time.Time{}, &time.ParseError{
		Layout:     "RFC3339, YYYY-MM-DD, or relative (1h, 2d, 30m)",
		Value:      s,
		LayoutElem: "",
		ValueElem:  s,
		Message:    "invalid time format",
	}
}

var relativeTimeRegex = regexp.MustCompile(`^(\d+)(m|h|d)$`)

func parseRelativeTime(s string) (time.Time, bool) {
	matches := relativeTimeRegex.FindStringSubmatch(strings.TrimSpace(s))
	if matches == nil {
		return time.Time{}, false
	}

	value, _ := strconv.Atoi(matches[1])
	unit := matches[2]

	var duration time.Duration
	switch unit {
	case "m":
		duration = time.Duration(value) * time.Minute
	case "h":
		duration = time.Duration(value) * time.Hour
	case "d":
		duration = time.Duration(value) * 24 * time.Hour
	default:
		return time.Time{}, false
	}

	return time.Now().Add(-duration), true
}
