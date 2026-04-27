package report

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/driftwatch/internal/drift"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes drift results to an output destination.
type Reporter struct {
	w      io.Writer
	format Format
}

// New creates a Reporter writing to w in the given format.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w, format: format}
}

// Write outputs the drift results according to the reporter's format.
func (r *Reporter) Write(results []drift.Result) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(results)
	default:
		return r.writeText(results)
	}
}

func (r *Reporter) writeText(results []drift.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(r.w, "No drift detected.")
		return nil
	}

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tRESOURCE ID\tDETAILS")
	fmt.Fprintln(tw, "------\t-----------\t-------")

	for _, res := range results {
		status := statusLabel(res.Status)
		details := res.Details
		if details == "" {
			details = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", status, res.ResourceID, details)
	}

	return tw.Flush()
}

func (r *Reporter) writeJSON(results []drift.Result) error {
	// Minimal JSON serialisation without importing encoding/json for brevity.
	fmt.Fprintln(r.w, "[")
	for i, res := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(r.w, "  {\"resource_id\": %q, \"status\": %q, \"details\": %q}%s\n",
			res.ResourceID, res.Status, res.Details, comma)
	}
	fmt.Fprintln(r.w, "]")
	return nil
}

func statusLabel(s drift.Status) string {
	switch s {
	case drift.StatusOK:
		return "OK"
	case drift.StatusMissing:
		return "MISSING"
	case drift.StatusDrifted:
		return "DRIFTED"
	default:
		return "UNKNOWN"
	}
}
