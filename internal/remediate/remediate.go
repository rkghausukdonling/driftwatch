// Package remediate provides suggestion generation for detected drift.
package remediate

import (
	"fmt"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Suggestion holds a human-readable remediation hint for a drifted resource.
type Suggestion struct {
	ResourceID string
	ResourceType string
	Status drift.Status
	Hint string
}

// Generate produces remediation suggestions for the given drift results.
func Generate(results []drift.Result) []Suggestion {
	suggestions := make([]Suggestion, 0, len(results))
	for _, r := range results {
		if r.Status == drift.StatusOK {
			continue
		}
		suggestions = append(suggestions, Suggestion{
			ResourceID:   r.ID,
			ResourceType: r.Type,
			Status:       r.Status,
			Hint:         hint(r),
		})
	}
	return suggestions
}

// hint builds a context-aware remediation message for a single result.
func hint(r drift.Result) string {
	switch r.Status {
	case drift.StatusMissing:
		return fmt.Sprintf(
			"Resource %q (%s) is declared in IaC but not found in the provider. "+
				"Run your IaC apply (e.g. `terraform apply`) to provision it, or remove the declaration.",
			r.ID, r.Type,
		)
	case drift.StatusDrifted:
		fields := driftedFields(r)
		if len(fields) > 0 {
			return fmt.Sprintf(
				"Resource %q (%s) has drifted fields: %s. "+
					"Reconcile by running `terraform apply` or updating your IaC definition.",
				r.ID, r.Type, strings.Join(fields, ", "),
			)
		}
		return fmt.Sprintf(
			"Resource %q (%s) has drifted from its IaC definition. Run `terraform apply` to reconcile.",
			r.ID, r.Type,
		)
	default:
		return fmt.Sprintf("Review resource %q (%s) manually.", r.ID, r.Type)
	}
}

// driftedFields returns the keys that differ between desired and actual state.
func driftedFields(r drift.Result) []string {
	var fields []string
	for k, desired := range r.Desired {
		actual, ok := r.Actual[k]
		if !ok || actual != desired {
			fields = append(fields, k)
		}
	}
	return fields
}
