// Package filter provides utilities for filtering drift detection results
// based on resource type, status, or ID patterns.
package filter

import (
	"strings"

	"github.com/your-org/driftwatch/internal/drift"
)

// Options holds the criteria used to filter drift results.
type Options struct {
	// Types restricts results to the given resource types (e.g. "aws_s3_bucket").
	// An empty slice means no type filtering.
	Types []string

	// Statuses restricts results to the given drift statuses.
	// An empty slice means no status filtering.
	Statuses []drift.Status

	// IDPrefix keeps only results whose resource ID starts with the given prefix.
	// An empty string means no prefix filtering.
	IDPrefix string
}

// Apply returns a new slice containing only the results that match all
// non-empty criteria in opts. The original slice is never modified.
func Apply(results []drift.Result, opts Options) []drift.Result {
	filtered := make([]drift.Result, 0, len(results))

	for _, r := range results {
		if !matchesType(r, opts.Types) {
			continue
		}
		if !matchesStatus(r, opts.Statuses) {
			continue
		}
		if !matchesIDPrefix(r, opts.IDPrefix) {
			continue
		}
		filtered = append(filtered, r)
	}

	return filtered
}

func matchesType(r drift.Result, types []string) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if strings.EqualFold(r.ResourceType, t) {
			return true
		}
	}
	return false
}

func matchesStatus(r drift.Result, statuses []drift.Status) bool {
	if len(statuses) == 0 {
		return true
	}
	for _, s := range statuses {
		if r.Status == s {
			return true
		}
	}
	return false
}

func matchesIDPrefix(r drift.Result, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(r.ResourceID, prefix)
}
