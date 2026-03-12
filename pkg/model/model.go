// Package model defines the core data structures for agent-rss.
package model

import "time"

// Feed represents a subscription entry.
type Feed struct {
	Name string `json:"name"`
	Src  string `json:"src"`
}

// Item represents a single RSS/Atom item with associated feed metadata.
type Item struct {
	Name    string    `json:"name"`
	Src     string    `json:"src"`
	Time    time.Time `json:"time"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Link    string    `json:"link"`
	ID      string    `json:"id"`
}
