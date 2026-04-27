package drift

import (
	"fmt"

	"github.com/driftwatch/internal/provider"
)

// Status represents the drift status of a single resource.
type Status int

const (
	StatusMatch   Status = iota // Resource matches IaC definition
	StatusDrifted               // Resource exists but has drifted
	StatusMissing               // Resource defined in IaC but not found
	StatusOrphan                // Resource found but not in IaC
)

func (s Status) String() string {
	switch s {
	case StatusMatch:
		return "match"
	case StatusDrifted:
		return "drifted"
	case StatusMissing:
		return "missing"
	case StatusOrphan:
		return "orphan"
	default:
		return "unknown"
	}
}

// Result holds the drift detection outcome for a single resource.
type Result struct {
	ResourceID string
	Status     Status
	Details    string
}

// Detector compares desired state (IaC definitions) against live provider state.
type Detector struct {
	provider provider.Provider
}

// New creates a Detector backed by the given provider.
func New(p provider.Provider) *Detector {
	return &Detector{provider: p}
}

// Detect checks each desired resource ID against live state and returns results.
func (d *Detector) Detect(desiredIDs []string) ([]Result, error) {
	if len(desiredIDs) == 0 {
		return nil, fmt.Errorf("detect: no resource IDs provided")
	}

	results := make([]Result, 0, len(desiredIDs))

	for _, id := range desiredIDs {
		res, err := d.provider.FetchState(id)
		if err != nil {
			results = append(results, Result{
				ResourceID: id,
				Status:     StatusMissing,
				Details:    fmt.Sprintf("provider error: %v", err),
			})
			continue
		}

		if res == nil {
			results = append(results, Result{
				ResourceID: id,
				Status:     StatusMissing,
				Details:    "resource not found in provider",
			})
			continue
		}

		results = append(results, Result{
			ResourceID: id,
			Status:     StatusMatch,
			Details:    fmt.Sprintf("state: %s", res.State),
		})
	}

	return results, nil
}
