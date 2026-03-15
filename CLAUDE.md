# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

agent-rss is a CLI tool for fetching and filtering RSS/Atom feeds, designed for AI agents. It supports subscription management, concurrent feed fetching, time-based and keyword filtering, and outputs structured JSON/NDJSON.

## Build and Development

### Build the binary
```bash
go build -o agent-rss ./cmd/agent-rss
```

### Run tests
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/storage
go test ./internal/filter
go test ./internal/output

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -v ./internal/filter -run TestParseTime
```

### Test the binary locally
```bash
# Build and test
go build -o agent-rss ./cmd/agent-rss
./agent-rss --help

# Test with custom feeds file
./agent-rss --file /tmp/test-feeds.txt list
```

### Code formatting and modernization

**IMPORTANT**: After modifying any Go code, always run these commands before committing:

```bash
# Format all Go code according to Go standards
go fmt ./...

# Modernize code to use latest Go idioms and fix deprecated patterns
go fix ./...
```

These tools ensure:
- **go fmt**: Consistent code formatting across the project
- **go fix**: Code uses modern Go patterns and APIs, deprecated code is updated

Run them automatically before every commit to maintain code quality.

## Architecture

### Package Structure

The codebase follows Go standard project layout:

- **cmd/agent-rss**: Main entry point, minimal logic
- **internal/**: Private application code
  - **cli/**: CLI command definitions using urfave/cli/v3
  - **storage/**: Feed subscription file management (plain text format)
  - **rss/**: RSS/Atom fetching using gofeed
  - **filter/**: Time and keyword filtering logic
  - **output/**: JSON/NDJSON serialization
  - **config/**: Configuration and path management
- **pkg/model**: Public data models (Feed, Item)

### Data Flow

1. **CLI Layer** (`internal/cli`): Parses commands and flags
2. **Storage Layer** (`internal/storage`): Reads feed subscriptions from text file
3. **Fetch Layer** (`internal/rss`): Fetches RSS/Atom feeds **concurrently**
4. **Filter Layer** (`internal/filter`): Filters items by time/keywords
5. **Output Layer** (`internal/output`): Serializes to JSON/NDJSON

### Key Design Decisions

**Concurrent Fetching** (cli.go:257-281):
- Multiple feeds are fetched in parallel using goroutines
- Each goroutine creates its own `rss.Fetcher` instance to avoid race conditions
- Results are synchronized using `sync.Mutex` and `sync.WaitGroup`
- This significantly improves performance when fetching multiple feeds

**Text-based Storage** (storage.go):
- Feeds stored in `~/.config/agent-rss/feeds.txt` as plain text
- Format: `name url` (space-separated, one per line)
- Comments start with `#`
- Feed names cannot contain whitespace

**Time Parsing Flexibility** (filter.go:65-119):
- Supports three formats:
  - RFC3339: `2026-03-12T08:30:00Z`
  - Date: `2026-03-12`
  - Relative: `1h`, `2d`, `30m` (hours, days, minutes ago)
- Relative times are calculated from `time.Now()`

**Output Formats**:
- NDJSON (default): One JSON object per line, suitable for streaming
- JSON: Single array, better for small result sets

## Release Process

This project uses automated releases via GitHub Actions:

1. **Create and push a git tag**:
   ```bash
   git tag -a v0.x.x -m "Release v0.x.x"
   git push origin v0.x.x
   ```

2. **GitHub Actions automatically**:
   - Builds binaries for 5 platforms (darwin-arm64, darwin-x64, linux-arm64, linux-x64, win32-x64)
   - Publishes platform-specific npm packages: `@atopos31/agent-rss-<platform>`
   - Publishes main npm package: `@atopos31/agent-rss`
   - Creates GitHub Release with auto-generated release notes

3. **Version must follow semver** (v0.1.4, v1.0.0, etc.)

## Important Implementation Details

### Concurrent Safety
When modifying `fetchCmd()` in `internal/cli/cli.go`, maintain concurrent safety:
- Create separate fetcher instances per goroutine
- Protect shared state with mutexes
- Use WaitGroup to synchronize completion

### Storage Format Constraints
The storage format imposes constraints:
- Feed names must not contain spaces (validated in storage.go:218-226)
- URLs are not validated (any string accepted)
- File format is fragile (space-delimited)

### Filter Behavior
Filters are AND-combined:
- Multiple `--title` flags are OR within titles
- Multiple `--content` flags are OR within content
- Time filters (`--since`, `--until`) are AND-combined
- All filter types are AND-combined together

### Output Size Limits
The tool is designed for AI agents which may have output size limits:
- Recommend writing output to file: `agent-rss fetch --all > /tmp/output.json`
- Then reading the file with agent's Read tool
- This pattern is documented in README and SKILL.md

## Feeds File Location

Default: `~/.config/agent-rss/feeds.txt`

Respects `XDG_CONFIG_HOME` environment variable:
```bash
# Custom location
XDG_CONFIG_HOME=/custom/path agent-rss list

# Or use --file flag
agent-rss --file /custom/feeds.txt list
```

## Dependencies

- **github.com/urfave/cli/v3**: CLI framework
- **github.com/mmcdole/gofeed**: RSS/Atom parsing
- Go 1.23+

## Skill Integration

This tool can be installed as a skill for AI agents:
```bash
npx skills add atopos31/agent-rss
```

The skill definition is in `skill/SKILL.md` and provides structured guidance for AI agents on how to use the tool effectively.
