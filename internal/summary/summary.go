// Package summary provides aggregate statistics over a set of drift detection results.
package summary

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Stats holds aggregate counts derived from a slice of DetectResult.
type Stats struct {
	Total   int
	Drifted int
	Missing int
	InSync  int
}

// Compute derives Stats from the provided results.
func Compute(results []drift.DetectResult) Stats {
	s := Stats{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case drift.StatusDrifted:
			s.Drifted++
		case drift.StatusMissing:
			s.Missing++
		case drift.StatusInSync:
			s.InSync++
		}
	}
	return s
}

// Write renders a human-readable summary table to w.
func Write(w io.Writer, s Stats) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SUMMARY")
	fmt.Fprintln(tw, "-------")
	fmt.Fprintf(tw, "Total\t%d\n", s.Total)
	fmt.Fprintf(tw, "In Sync\t%d\n", s.InSync)
	fmt.Fprintf(tw, "Drifted\t%d\n", s.Drifted)
	fmt.Fprintf(tw, "Missing\t%d\n", s.Missing)
	return tw.Flush()
}

// ExitCode returns a non-zero exit code when drift or missing resources exist.
func ExitCode(s Stats) int {
	if s.Drifted > 0 || s.Missing > 0 {
		return 1
	}
	return 0
}
