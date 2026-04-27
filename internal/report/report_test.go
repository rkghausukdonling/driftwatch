package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "vpc-abc123", Status: drift.StatusOK, Details: ""},
		{ResourceID: "sg-def456", Status: drift.StatusMissing, Details: "not found in provider"},
		{ResourceID: "subnet-ghi789", Status: drift.StatusDrifted, Details: "tag mismatch"},
	}
}

func TestWriteText_NoResults(t *testing.T) {
	var buf bytes.Buffer
	r := report.New(&buf, report.FormatText)
	if err := r.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected 'No drift detected', got: %q", buf.String())
	}
}

func TestWriteText_WithResults(t *testing.T) {
	var buf bytes.Buffer
	r := report.New(&buf, report.FormatText)
	if err := r.Write(makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"OK", "MISSING", "DRIFTED", "vpc-abc123", "sg-def456", "subnet-ghi789"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestWriteText_EmptyDetails(t *testing.T) {
	var buf bytes.Buffer
	r := report.New(&buf, report.FormatText)
	results := []drift.Result{{ResourceID: "igw-001", Status: drift.StatusOK, Details: ""}}
	if err := r.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "-") {
		t.Errorf("expected placeholder '-' for empty details")
	}
}

func TestWriteJSON_WithResults(t *testing.T) {
	var buf bytes.Buffer
	r := report.New(&buf, report.FormatJSON)
	if err := r.Write(makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"resource_id", "status", "details", "vpc-abc123", "MISSING", "DRIFTED"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected JSON output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestNew_NilWriter(t *testing.T) {
	// Should not panic when nil writer is provided.
	r := report.New(nil, report.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
