// Package rss provides RSS/Atom feed fetching and parsing.
package rss

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/atopos31/agent-rss/pkg/model"
)

// Fetcher handles RSS/Atom feed retrieval.
type Fetcher struct {
	parser *gofeed.Parser
}

// New creates a new Fetcher.
func New() *Fetcher {
	return &Fetcher{
		parser: gofeed.NewParser(),
	}
}

// Fetch retrieves and parses an RSS/Atom feed from the given URL.
func (f *Fetcher) Fetch(ctx context.Context, feed model.Feed) ([]model.Item, error) {
	parsed, err := f.parser.ParseURLWithContext(feed.Src, ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", feed.Name, err)
	}

	items := make([]model.Item, 0, len(parsed.Items))

	for _, item := range parsed.Items {
		items = append(items, convertItem(feed, item))
	}

	return items, nil
}

func convertItem(feed model.Feed, item *gofeed.Item) model.Item {
	t := time.Time{}
	if item.PublishedParsed != nil {
		t = item.PublishedParsed.Local()
	} else if item.UpdatedParsed != nil {
		t = item.UpdatedParsed.Local()
	}

	content := item.Description
	if item.Content != "" {
		content = item.Content
	}

	id := item.GUID
	if id == "" {
		id = item.Link
	}

	return model.Item{
		Name:    feed.Name,
		Src:     feed.Src,
		Time:    t,
		Title:   item.Title,
		Content: content,
		Link:    item.Link,
		ID:      id,
	}
}
