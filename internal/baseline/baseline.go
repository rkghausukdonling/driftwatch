// Package baseline provides functionality to save and compare drift scan
// results against a known-good baseline snapshot.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/driftwatch/internal/drift"
)

// Snapshot represents a saved baseline of drift detection results.
type Snapshot struct {
	CreatedAt time.Time            `json:"created_at"`
	Provider  string               `json:"provider"`
	Results   []drift.Result       `json:"results"`
}

// Save writes the given results to the specified file as a JSON baseline snapshot.
func Save(path, provider string, results []drift.Result) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Provider:  provider,
		Results:   results,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write file %q: %w", path, err)
	}
	return nil
}

// Load reads a baseline snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: file not found: %q", path)
		}
		return nil, fmt.Errorf("baseline: read file %q: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal snapshot: %w", err)
	}
	return &snap, nil
}

// Compare returns the results from current that differ from the baseline snapshot.
// A result is considered new or changed if its ID or Status does not match the baseline.
func Compare(snap *Snapshot, current []drift.Result) []drift.Result {
	index := make(map[string]drift.Result, len(snap.Results))
	for _, r := range snap.Results {
		index[r.ID] = r
	}
	var diffs []drift.Result
	for _, r := range current {
		base, found := index[r.ID]
		if !found || base.Status != r.Status {
			diffs = append(diffs, r)
		}
	}
	return diffs
}
