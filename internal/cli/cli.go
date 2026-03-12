// Package cli provides the command-line interface using urfave/cli/v3.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/atopos31/agent-rss/internal/config"
	"github.com/atopos31/agent-rss/internal/filter"
	"github.com/atopos31/agent-rss/internal/output"
	"github.com/atopos31/agent-rss/internal/rss"
	"github.com/atopos31/agent-rss/internal/storage"
	"github.com/atopos31/agent-rss/pkg/model"
)

// App creates and returns the CLI application.
func App() *cli.Command {
	return &cli.Command{
		Name:  "agent-rss",
		Usage: "CLI RSS tool for AI agents",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "path to feeds file",
				Value:   config.DefaultFeedsPath(),
			},
		},
		Commands: []*cli.Command{
			addCmd(),
			listCmd(),
			getCmd(),
			updateCmd(),
			removeCmd(),
			fetchCmd(),
		},
	}
}

func getStore(cmd *cli.Command) *storage.Store {
	return storage.New(cmd.String("file"))
}

func addCmd() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new feed subscription",
		ArgsUsage: "<name> <src>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := cmd.Args()
			if args.Len() < 2 {
				return fmt.Errorf("usage: agent-rss add <name> <src>")
			}

			name := args.Get(0)
			src := args.Get(1)

			store := getStore(cmd)
			if err := store.Add(name, src); err != nil {
				return err
			}

			fmt.Printf("Added feed: %s\n", name)
			return nil
		},
	}
}

func listCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all feed subscriptions",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			store := getStore(cmd)
			feeds, err := store.List()
			if err != nil {
				return err
			}

			if len(feeds) == 0 {
				fmt.Println("No feeds configured")
				return nil
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(feeds)
		},
	}
}

func getCmd() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Get a feed by name",
		ArgsUsage: "<name>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := cmd.Args()
			if args.Len() < 1 {
				return fmt.Errorf("usage: agent-rss get <name>")
			}

			name := args.Get(0)
			store := getStore(cmd)
			feed, err := store.Get(name)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(feed)
		},
	}
}

func updateCmd() *cli.Command {
	return &cli.Command{
		Name:      "update",
		Usage:     "Update an existing feed",
		ArgsUsage: "<name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "new name for the feed",
			},
			&cli.StringFlag{
				Name:  "src",
				Usage: "new source URL",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := cmd.Args()
			if args.Len() < 1 {
				return fmt.Errorf("usage: agent-rss update <name> [--name <new>] [--src <newsrc>]")
			}

			name := args.Get(0)
			store := getStore(cmd)

			var newName, newSrc *string
			if cmd.IsSet("name") {
				n := cmd.String("name")
				newName = &n
			}
			if cmd.IsSet("src") {
				s := cmd.String("src")
				newSrc = &s
			}

			if newName == nil && newSrc == nil {
				return fmt.Errorf("specify at least one of --name or --src")
			}

			if err := store.Update(name, newName, newSrc); err != nil {
				return err
			}

			fmt.Printf("Updated feed: %s\n", name)
			return nil
		},
	}
}

func removeCmd() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "Remove a feed subscription",
		ArgsUsage: "<name>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := cmd.Args()
			if args.Len() < 1 {
				return fmt.Errorf("usage: agent-rss remove <name>")
			}

			name := args.Get(0)
			store := getStore(cmd)
			if err := store.Remove(name); err != nil {
				return err
			}

			fmt.Printf("Removed feed: %s\n", name)
			return nil
		},
	}
}

func fetchCmd() *cli.Command {
	return &cli.Command{
		Name:  "fetch",
		Usage: "Fetch RSS items",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "fetch specific feed by name",
			},
			&cli.BoolFlag{
				Name:  "all",
				Usage: "fetch all feeds",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "output format: json or ndjson",
				Value: "ndjson",
			},
			&cli.StringFlag{
				Name:  "since",
				Usage: "filter items published after this time (RFC3339 or YYYY-MM-DD)",
			},
			&cli.StringFlag{
				Name:  "until",
				Usage: "filter items published before this time (RFC3339 or YYYY-MM-DD)",
			},
			&cli.StringSliceFlag{
				Name:  "title",
				Usage: "filter by title keyword (can be specified multiple times)",
			},
			&cli.StringSliceFlag{
				Name:  "content",
				Usage: "filter by content keyword (can be specified multiple times)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			store := getStore(cmd)
			fetcher := rss.New()

			var feeds []model.Feed

			if cmd.Bool("all") {
				var err error
				feeds, err = store.List()
				if err != nil {
					return err
				}
				if len(feeds) == 0 {
					return fmt.Errorf("no feeds configured")
				}
			} else if cmd.IsSet("name") {
				feed, err := store.Get(cmd.String("name"))
				if err != nil {
					return err
				}
				feeds = []model.Feed{feed}
			} else {
				return fmt.Errorf("specify --name <name> or --all")
			}

			opts, err := buildFilterOptions(cmd)
			if err != nil {
				return err
			}

			var allItems []model.Item
			for _, feed := range feeds {
				items, err := fetcher.Fetch(ctx, feed)
				if err != nil {
					fmt.Fprintf(os.Stderr, "warning: %v\n", err)
					continue
				}
				allItems = append(allItems, items...)
			}

			filtered := filter.Filter(allItems, opts)

			format := output.ParseFormat(cmd.String("format"))
			writer := output.New(os.Stdout, format)
			return writer.Write(filtered)
		},
	}
}

func buildFilterOptions(cmd *cli.Command) (filter.Options, error) {
	var opts filter.Options

	if cmd.IsSet("since") {
		t, err := filter.ParseTime(cmd.String("since"))
		if err != nil {
			return opts, fmt.Errorf("invalid --since: %w", err)
		}
		opts.Since = &t
	}

	if cmd.IsSet("until") {
		t, err := filter.ParseTime(cmd.String("until"))
		if err != nil {
			return opts, fmt.Errorf("invalid --until: %w", err)
		}
		opts.Until = &t
	}

	opts.Titles = cmd.StringSlice("title")
	opts.Contents = cmd.StringSlice("content")

	return opts, nil
}
