// Package filter provides time and keyword filtering for RSS items.
package filter

import (
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

// ParseTime attempts to parse a time string in RFC3339 or YYYY-MM-DD format.
func ParseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		return t, nil
	}

	return time.Time{}, &time.ParseError{
		Layout:     "RFC3339 or YYYY-MM-DD",
		Value:      s,
		LayoutElem: "",
		ValueElem:  s,
		Message:    "invalid time format",
	}
}
