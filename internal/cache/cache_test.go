package cache_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/cache"
)

func newTempCache(t *testing.T) *cache.Cache {
	t.Helper()
	dir := t.TempDir()
	c, err := cache.New(dir)
	if err != nil {
		t.Fatalf("cache.New: %v", err)
	}
	return c
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := t.TempDir() + "/subdir/cache"
	_, err := cache.New(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("expected directory %q to be created", dir)
	}
}

func TestSetAndGet_RoundTrip(t *testing.T) {
	c := newTempCache(t)

	entry := cache.Entry{
		ResourceID: "i-abc123",
		Provider:   "aws",
		Attributes: map[string]string{"instance_type": "t3.micro", "region": "us-east-1"},
		CachedAt:   time.Now().UTC().Truncate(time.Second),
	}

	if err := c.Set(entry); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, found, err := c.Get("aws", "i-abc123")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !found {
		t.Fatal("expected entry to be found")
	}
	if got.ResourceID != entry.ResourceID {
		t.Errorf("ResourceID: got %q, want %q", got.ResourceID, entry.ResourceID)
	}
	if got.Attributes["instance_type"] != "t3.micro" {
		t.Errorf("Attributes[instance_type]: got %q, want %q", got.Attributes["instance_type"], "t3.micro")
	}
}

func TestGet_NotFound(t *testing.T) {
	c := newTempCache(t)

	_, found, err := c.Get("aws", "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected entry not to be found")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c := newTempCache(t)

	entry := cache.Entry{ResourceID: "sg-001", Provider: "aws", Attributes: map[string]string{}}
	if err := c.Set(entry); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := c.Delete("aws", "sg-001"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, found, err := c.Get("aws", "sg-001")
	if err != nil {
		t.Fatalf("Get after delete: %v", err)
	}
	if found {
		t.Fatal("expected entry to be absent after delete")
	}
}

func TestDelete_NonExistent_NoError(t *testing.T) {
	c := newTempCache(t)
	if err := c.Delete("aws", "does-not-exist"); err != nil {
		t.Fatalf("expected no error deleting nonexistent entry, got %v", err)
	}
}
