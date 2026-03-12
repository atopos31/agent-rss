package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStore_CRUD(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feeds.txt")
	store := New(path)

	// Test Add
	if err := store.Add("test1", "https://example.com/rss"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Test List
	feeds, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}
	if feeds[0].Name != "test1" || feeds[0].Src != "https://example.com/rss" {
		t.Fatalf("unexpected feed: %+v", feeds[0])
	}

	// Test Get
	feed, err := store.Get("test1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if feed.Name != "test1" {
		t.Fatalf("unexpected feed name: %s", feed.Name)
	}

	// Test Update
	newName := "test1-updated"
	if err := store.Update("test1", &newName, nil); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	feed, err = store.Get("test1-updated")
	if err != nil {
		t.Fatalf("Get updated failed: %v", err)
	}
	if feed.Name != "test1-updated" {
		t.Fatalf("update did not change name")
	}

	// Test Remove
	if err := store.Remove("test1-updated"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	feeds, err = store.List()
	if err != nil {
		t.Fatalf("List after remove failed: %v", err)
	}
	if len(feeds) != 0 {
		t.Fatalf("expected 0 feeds after remove, got %d", len(feeds))
	}
}

func TestStore_AddDuplicate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feeds.txt")
	store := New(path)

	if err := store.Add("dup", "https://example.com/rss"); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}

	if err := store.Add("dup", "https://example.com/rss2"); err != ErrFeedExists {
		t.Fatalf("expected ErrFeedExists, got: %v", err)
	}
}

func TestStore_RemoveNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feeds.txt")
	store := New(path)

	if err := store.Remove("nonexistent"); err != ErrFeedNotFound {
		t.Fatalf("expected ErrFeedNotFound, got: %v", err)
	}
}

func TestStore_ParseComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feeds.txt")

	content := `# This is a comment
test1 https://example.com/rss

# Another comment
test2 https://example.com/atom
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	store := New(path)
	feeds, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(feeds) != 2 {
		t.Fatalf("expected 2 feeds, got %d", len(feeds))
	}
}

func TestStore_ValidationErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feeds.txt")
	store := New(path)

	tests := []struct {
		name    string
		src     string
		wantErr error
	}{
		{"", "https://example.com", ErrEmptyName},
		{"test", "", ErrEmptySrc},
		{"has space", "https://example.com", ErrNameHasSpace},
	}

	for _, tt := range tests {
		err := store.Add(tt.name, tt.src)
		if err != tt.wantErr {
			t.Errorf("Add(%q, %q): expected %v, got %v", tt.name, tt.src, tt.wantErr, err)
		}
	}
}
