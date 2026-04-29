// Package cache provides a simple file-backed cache for storing provider
// state snapshots to enable drift comparison across multiple runs.
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single cached resource state snapshot.
type Entry struct {
	ResourceID string            `json:"resource_id"`
	Provider   string            `json:"provider"`
	Attributes map[string]string `json:"attributes"`
	CachedAt   time.Time         `json:"cached_at"`
}

// Cache manages persisted state entries on disk.
type Cache struct {
	dir string
}

// New creates a Cache that stores entries under dir.
// The directory is created if it does not exist.
func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cache: create directory %q: %w", dir, err)
	}
	return &Cache{dir: dir}, nil
}

// Set writes an Entry to disk, keyed by provider and resource ID.
func (c *Cache) Set(e Entry) error {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal entry: %w", err)
	}
	if err := os.WriteFile(c.path(e.Provider, e.ResourceID), data, 0o644); err != nil {
		return fmt.Errorf("cache: write entry: %w", err)
	}
	return nil
}

// Get retrieves a previously cached Entry. Returns (zero, false, nil) when the
// entry does not exist, and (zero, false, err) on unexpected errors.
func (c *Cache) Get(provider, resourceID string) (Entry, bool, error) {
	data, err := os.ReadFile(c.path(provider, resourceID))
	if os.IsNotExist(err) {
		return Entry{}, false, nil
	}
	if err != nil {
		return Entry{}, false, fmt.Errorf("cache: read entry: %w", err)
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return Entry{}, false, fmt.Errorf("cache: unmarshal entry: %w", err)
	}
	return e, true, nil
}

// Delete removes a cached entry. It is not an error if the entry does not exist.
func (c *Cache) Delete(provider, resourceID string) error {
	err := os.Remove(c.path(provider, resourceID))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cache: delete entry: %w", err)
	}
	return nil
}

// path returns the file path for a given provider + resource pair.
func (c *Cache) path(provider, resourceID string) string {
	filename := fmt.Sprintf("%s__%s.json", provider, resourceID)
	return filepath.Join(c.dir, filename)
}
