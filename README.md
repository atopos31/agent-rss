# agent-rss

CLI RSS tool for AI agents. Subscribe to feeds, fetch RSS/Atom content, and filter by time or keywords.

## Installation

### npm

```bash
npm install -g @atopos31/agent-rss
```

### Go

```bash
go install github.com/atopos31/agent-rss/cmd/agent-rss@latest
```

### From Source

```bash
git clone https://github.com/atopos31/agent-rss.git
cd agent-rss
go build -o agent-rss ./cmd/agent-rss
```

## Usage

### Subscription Management

```bash
# Add a feed
agent-rss add hn https://news.ycombinator.com/rss

# List all feeds
agent-rss list

# Get a specific feed
agent-rss get hn

# Update a feed
agent-rss update hn --src https://news.ycombinator.com/rss

# Remove a feed
agent-rss remove hn
```

### Fetching RSS

```bash
# Fetch a specific feed
agent-rss fetch --name hn

# Fetch all feeds
agent-rss fetch --all

# Output as JSON array (default is NDJSON)
agent-rss fetch --all --format json
```

### Filtering

```bash
# Filter by time range
agent-rss fetch --all --since 2026-03-12
agent-rss fetch --all --since 2026-03-12T08:00:00Z --until 2026-03-12T18:00:00Z

# Filter by title keyword
agent-rss fetch --all --title "AI"

# Filter by content keyword
agent-rss fetch --all --content "machine learning"

# Combine filters
agent-rss fetch --all --since 2026-03-12 --title "Go" --title "Rust"
```

### Global Options

```bash
# Use a custom feeds file
agent-rss --file /path/to/feeds.txt list
```

## Feeds File Format

Feeds are stored in `~/.config/agent-rss/feeds.txt`:

```
# Comments start with #
hn https://news.ycombinator.com/rss
golang https://blog.golang.org/feed.atom
```

## Output Format

### NDJSON (default)

```json
{"name":"hn","src":"https://...","time":"2026-03-12T15:30:00+08:00","title":"...","content":"...","link":"...","id":"..."}
{"name":"hn","src":"https://...","time":"2026-03-12T14:20:00+08:00","title":"...","content":"...","link":"...","id":"..."}
```

### JSON

```json
[
  {"name":"hn","src":"https://...","time":"2026-03-12T15:30:00+08:00","title":"...","content":"...","link":"...","id":"..."},
  {"name":"hn","src":"https://...","time":"2026-03-12T14:20:00+08:00","title":"...","content":"...","link":"...","id":"..."}
]
```

## Find RSS Feeds

Looking for RSS feeds to subscribe? Check out [awesome-rsshub-routes](https://github.com/JackyST0/awesome-rsshub-routes) for a curated list of RSS feeds across various categories.

## License

MIT
