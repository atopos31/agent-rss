// Package config provides configuration and path management.
package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultFileName = "feeds.txt"
	DefaultDirName  = "agent-rss"
)

// DefaultFeedsPath returns the default path for the feeds file.
// It uses $XDG_CONFIG_HOME/agent-rss/feeds.txt or ~/.config/agent-rss/feeds.txt.
func DefaultFeedsPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return DefaultFileName
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, DefaultDirName, DefaultFileName)
}

// EnsureDir creates the parent directory for the given path if it doesn't exist.
func EnsureDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}
