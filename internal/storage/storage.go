// Package storage provides feed subscription file management.
package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atopos31/agent-rss/internal/config"
	"github.com/atopos31/agent-rss/pkg/model"
)

var (
	ErrFeedNotFound  = errors.New("feed not found")
	ErrFeedExists    = errors.New("feed already exists")
	ErrInvalidFormat = errors.New("invalid feed format: expected 'name src'")
	ErrEmptyName     = errors.New("feed name cannot be empty")
	ErrEmptySrc      = errors.New("feed source cannot be empty")
	ErrNameHasSpace  = errors.New("feed name cannot contain whitespace")
)

// Store manages feed subscriptions in a text file.
type Store struct {
	path string
}

// New creates a new Store with the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Path returns the file path used by this store.
func (s *Store) Path() string {
	return s.path
}

// List returns all feeds from the subscription file.
func (s *Store) List() ([]model.Feed, error) {
	f, err := os.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open feeds file: %w", err)
	}
	defer f.Close()

	return parseFeeds(f)
}

// Get retrieves a single feed by name.
func (s *Store) Get(name string) (model.Feed, error) {
	feeds, err := s.List()
	if err != nil {
		return model.Feed{}, err
	}

	for _, feed := range feeds {
		if feed.Name == name {
			return feed, nil
		}
	}
	return model.Feed{}, ErrFeedNotFound
}

// Add adds a new feed subscription.
func (s *Store) Add(name, src string) error {
	if err := validateFeed(name, src); err != nil {
		return err
	}

	feeds, err := s.List()
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		if feed.Name == name {
			return ErrFeedExists
		}
	}

	feeds = append(feeds, model.Feed{Name: name, Src: src})
	return s.write(feeds)
}

// Update modifies an existing feed.
func (s *Store) Update(name string, newName, newSrc *string) error {
	feeds, err := s.List()
	if err != nil {
		return err
	}

	idx := -1
	for i, feed := range feeds {
		if feed.Name == name {
			idx = i
			break
		}
	}

	if idx == -1 {
		return ErrFeedNotFound
	}

	updated := feeds[idx]

	if newName != nil {
		if err := validateName(*newName); err != nil {
			return err
		}
		for i, feed := range feeds {
			if i != idx && feed.Name == *newName {
				return ErrFeedExists
			}
		}
		updated.Name = *newName
	}

	if newSrc != nil {
		if err := validateSrc(*newSrc); err != nil {
			return err
		}
		updated.Src = *newSrc
	}

	feeds[idx] = updated
	return s.write(feeds)
}

// Remove deletes a feed by name.
func (s *Store) Remove(name string) error {
	feeds, err := s.List()
	if err != nil {
		return err
	}

	idx := -1
	for i, feed := range feeds {
		if feed.Name == name {
			idx = i
			break
		}
	}

	if idx == -1 {
		return ErrFeedNotFound
	}

	feeds = append(feeds[:idx], feeds[idx+1:]...)
	return s.write(feeds)
}

func (s *Store) write(feeds []model.Feed) error {
	if err := config.EnsureDir(s.path); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	f, err := os.Create(s.path)
	if err != nil {
		return fmt.Errorf("create feeds file: %w", err)
	}
	defer f.Close()

	for _, feed := range feeds {
		if _, err := fmt.Fprintf(f, "%s %s\n", feed.Name, feed.Src); err != nil {
			return fmt.Errorf("write feed: %w", err)
		}
	}
	return nil
}

func parseFeeds(r io.Reader) ([]model.Feed, error) {
	var feeds []model.Feed
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: %w", lineNum, ErrInvalidFormat)
		}

		name := strings.TrimSpace(parts[0])
		src := strings.TrimSpace(parts[1])

		if name == "" || src == "" {
			return nil, fmt.Errorf("line %d: %w", lineNum, ErrInvalidFormat)
		}

		feeds = append(feeds, model.Feed{Name: name, Src: src})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read feeds file: %w", err)
	}

	return feeds, nil
}

func validateFeed(name, src string) error {
	if err := validateName(name); err != nil {
		return err
	}
	return validateSrc(src)
}

func validateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	if strings.ContainsAny(name, " \t\n\r") {
		return ErrNameHasSpace
	}
	return nil
}

func validateSrc(src string) error {
	if src == "" {
		return ErrEmptySrc
	}
	return nil
}
