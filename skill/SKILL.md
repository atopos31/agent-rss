---
name: rss-skill
description: Use this skill when the user wants to fetch, read, or monitor RSS/Atom feeds. This includes subscribing to RSS sources, fetching latest news/articles from feeds, filtering feed items by time (past hours/days) or keywords, and managing feed subscriptions. Use this skill when the user mentions RSS, feeds, news subscriptions, or wants to stay updated on websites/blogs.
---

# RSS Feed Management with agent-rss

## Overview

agent-rss is a CLI tool for fetching and filtering RSS/Atom feeds. It outputs structured JSON for easy processing.

## Installation Check

First, check if agent-rss is installed:

```bash
which agent-rss || echo "not installed"
```

If not installed:
```bash
npm install -g @atopos31/agent-rss
```

## Core Commands

### Managing Subscriptions

```bash
# Add a feed
agent-rss add <name> <url>

# List all feeds
agent-rss list

# Get a specific feed
agent-rss get <name>

# Update a feed
agent-rss update <name> --src <new-url>

# Remove a feed
agent-rss remove <name>
```

### Fetching RSS

```bash
# Fetch specific feed
agent-rss fetch --name <name>

# Fetch all feeds
agent-rss fetch --all

# Output as JSON array (default is NDJSON)
agent-rss fetch --all --format json
```

### Time Filtering

```bash
# Relative time (recommended for agents)
agent-rss fetch --all --since 1h      # past 1 hour
agent-rss fetch --all --since 2d      # past 2 days
agent-rss fetch --all --since 30m     # past 30 minutes

# Absolute time
agent-rss fetch --all --since 2026-03-12
agent-rss fetch --all --since 2026-03-12T08:00:00+08:00
```

### Keyword Filtering

```bash
# Filter by title
agent-rss fetch --all --title "AI" --title "ML"

# Filter by content
agent-rss fetch --all --content "machine learning"

# Combine filters
agent-rss fetch --all --since 1d --title "AI"
```

## Best Practice: Output to File

**IMPORTANT**: CLI output may be truncated due to size limits. Always write output to a file, then read it:

```bash
# Step 1: Write to file
agent-rss fetch --all --since 1d > /tmp/rss-output.json

# Step 2: Use Read tool to access full content
# Read /tmp/rss-output.json
```

This ensures no data is lost due to output truncation.

## Output Format

Each item contains:
```json
{
  "name": "feed-name",
  "src": "https://example.com/rss",
  "time": "2026-03-12T15:30:00+08:00",
  "title": "Article Title",
  "content": "Article content or summary",
  "link": "https://example.com/article",
  "id": "unique-id"
}
```

## Common Workflows

### Check for Recent News

```bash
agent-rss fetch --all --since 1h > /tmp/recent.json
# Then read /tmp/recent.json
```

### Search for Specific Topics

```bash
agent-rss fetch --all --since 1d --title "AI" --title "LLM" > /tmp/ai-news.json
# Then read /tmp/ai-news.json
```

### Daily News Summary

```bash
agent-rss fetch --all --since 24h --format json > /tmp/daily.json
# Then read /tmp/daily.json and summarize
```

## Finding RSS Feeds

For curated RSS feeds, see: https://github.com/JackyST0/awesome-rsshub-routes

## Troubleshooting

- If a feed times out, try fetching it individually with `--name`
- Use `--format json` for JSON array output
- Times are displayed in local timezone
