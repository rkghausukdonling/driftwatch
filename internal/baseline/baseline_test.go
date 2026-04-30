package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/baseline"
	"github.com/yourusername/driftwatch/internal/drift"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ID: "res-1", Type: "aws_instance", Status: drift.StatusMatch},
		{ID: "res-2", Type: "aws_s3_bucket", Status: drift.StatusDrifted},
		{ID: "res-3", Type: "aws_vpc", Status: drift.StatusMissing},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	results := makeResults()

	if err := baseline.Save(path, "aws", results); err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: unexpected error: %v", err)
	}
	if snap.Provider != "aws" {
		t.Errorf("provider: got %q, want %q", snap.Provider, "aws")
	}
	if len(snap.Results) != len(results) {
		t.Errorf("results count: got %d, want %d", len(snap.Results), len(results))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)

	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestCompare_NoDiff(t *testing.T) {
	results := makeResults()
	snap := &baseline.Snapshot{CreatedAt: time.Now(), Provider: "aws", Results: results}
	diffs := baseline.Compare(snap, results)
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(diffs))
	}
}

func TestCompare_StatusChanged(t *testing.T) {
	original := makeResults()
	snap := &baseline.Snapshot{CreatedAt: time.Now(), Provider: "aws", Results: original}

	current := makeResults()
	current[0].Status = drift.StatusDrifted // was StatusMatch

	diffs := baseline.Compare(snap, current)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].ID != "res-1" {
		t.Errorf("expected diff for res-1, got %q", diffs[0].ID)
	}
}

func TestCompare_NewResource(t *testing.T) {
	original := makeResults()
	snap := &baseline.Snapshot{CreatedAt: time.Now(), Provider: "aws", Results: original}

	current := append(makeResults(), drift.Result{ID: "res-99", Type: "aws_lambda", Status: drift.StatusDrifted})

	diffs := baseline.Compare(snap, current)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff (new resource), got %d", len(diffs))
	}
	if diffs[0].ID != "res-99" {
		t.Errorf("expected diff for res-99, got %q", diffs[0].ID)
	}
}
