package drift

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/provider"
)

// Status represents the drift state of a single resource.
type Status string

const (
	StatusOK      Status = "ok"
	StatusMissing Status = "missing"
	StatusDrifted Status = "drifted"
)

// Result holds the outcome of checking one resource for drift.
type Result struct {
	ResourceID string
	Status     Status
	Details    string
}

// Detector checks resources against a provider for drift.
type Detector struct {
	p provider.Provider
}

// New creates a Detector backed by the given provider.
func New(p provider.Provider) (*Detector, error) {
	if p == nil {
		return nil, fmt.Errorf("provider must not be nil")
	}
	return &Detector{p: p}, nil
}

// Detect checks each resource ID and returns a Result per resource.
func (d *Detector) Detect(ctx context.Context, resourceIDs []string) ([]Result, error) {
	if len(resourceIDs) == 0 {
		return nil, nil
	}

	results := make([]Result, 0, len(resourceIDs))

	for _, id := range resourceIDs {
		res, err := d.p.FetchState(ctx, id)
		if err != nil {
			results = append(results, Result{
				ResourceID: id,
				Status:     StatusMissing,
				Details:    fmt.Sprintf("fetch error: %v", err),
			})
			continue
		}

		if res == nil {
			results = append(results, Result{
				ResourceID: id,
				Status:     StatusMissing,
				Details:    "not found in provider",
			})
			continue
		}

		if res.Drifted {
			results = append(results, Result{
				ResourceID: id,
				Status:     StatusDrifted,
				Details:    res.DriftDetails,
			})
			continue
		}

		results = append(results, Result{
			ResourceID: id,
			Status:     StatusOK,
		})
	}

	return results, nil
}
