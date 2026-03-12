// Command agent-rss is a CLI RSS tool for AI agents.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/atopos31/agent-rss/internal/cli"
)

func main() {
	if err := cli.App().Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
