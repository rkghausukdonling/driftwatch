package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/summary"
)

func makeResults(statuses ...drift.Status) []drift.DetectResult {
	results := make([]drift.DetectResult, len(statuses))
	for i, s := range statuses {
		results[i] = drift.DetectResult{ResourceID: fmt.Sprintf("res-%d", i), Status: s}
	}
	return results
}

func TestCompute_Empty(t *testing.T) {
	s := summary.Compute(nil)
	if s.Total != 0 || s.InSync != 0 || s.Drifted != 0 || s.Missing != 0 {
		t.Errorf("expected all zeros, got %+v", s)
	}
}

func TestCompute_Mixed(t *testing.T) {
	results := []drift.DetectResult{
		{ResourceID: "a", Status: drift.StatusInSync},
		{ResourceID: "b", Status: drift.StatusDrifted},
		{ResourceID: "c", Status: drift.StatusMissing},
		{ResourceID: "d", Status: drift.StatusInSync},
	}
	s := summary.Compute(results)
	if s.Total != 4 {
		t.Errorf("Total: want 4, got %d", s.Total)
	}
	if s.InSync != 2 {
		t.Errorf("InSync: want 2, got %d", s.InSync)
	}
	if s.Drifted != 1 {
		t.Errorf("Drifted: want 1, got %d", s.Drifted)
	}
	if s.Missing != 1 {
		t.Errorf("Missing: want 1, got %d", s.Missing)
	}
}

func TestWrite_ContainsFields(t *testing.T) {
	s := summary.Stats{Total: 3, InSync: 1, Drifted: 1, Missing: 1}
	var buf bytes.Buffer
	if err := summary.Write(&buf, s); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"SUMMARY", "Total", "In Sync", "Drifted", "Missing"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got:\n%s", want, out)
		}
	}
}

func TestExitCode(t *testing.T) {
	if code := summary.ExitCode(summary.Stats{InSync: 5}); code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
	if code := summary.ExitCode(summary.Stats{Drifted: 1}); code != 1 {
		t.Errorf("expected 1, got %d", code)
	}
	if code := summary.ExitCode(summary.Stats{Missing: 1}); code != 1 {
		t.Errorf("expected 1, got %d", code)
	}
}
