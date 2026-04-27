package drift_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/provider"
	_ "github.com/driftwatch/internal/provider/mock"
)

func newMockProvider(t *testing.T) provider.Provider {
	t.Helper()
	p, err := provider.New("mock", nil)
	if err != nil {
		t.Fatalf("failed to create mock provider: %v", err)
	}
	return p
}

func TestDetect_EmptyIDs(t *testing.T) {
	d := drift.New(newMockProvider(t))
	_, err := d.Detect([]string{})
	if err == nil {
		t.Fatal("expected error for empty IDs, got nil")
	}
}

func TestDetect_FoundResource(t *testing.T) {
	d := drift.New(newMockProvider(t))
	// "instance-1" is a default resource in the mock provider
	results, err := d.Detect([]string{"instance-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != drift.StatusMatch {
		t.Errorf("expected StatusMatch, got %s", results[0].Status)
	}
}

func TestDetect_MissingResource(t *testing.T) {
	d := drift.New(newMockProvider(t))
	results, err := d.Detect([]string{"nonexistent-999"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != drift.StatusMissing {
		t.Errorf("expected StatusMissing, got %s", results[0].Status)
	}
}

func TestDetect_MixedResources(t *testing.T) {
	d := drift.New(newMockProvider(t))
	results, err := d.Detect([]string{"instance-1", "ghost-resource"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	statuses := map[string]drift.Status{}
	for _, r := range results {
		statuses[r.ResourceID] = r.Status
	}
	if statuses["instance-1"] != drift.StatusMatch {
		t.Errorf("instance-1: expected match, got %s", statuses["instance-1"])
	}
	if statuses["ghost-resource"] != drift.StatusMissing {
		t.Errorf("ghost-resource: expected missing, got %s", statuses["ghost-resource"])
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		status drift.Status
		want   string
	}{
		{drift.StatusMatch, "match"},
		{drift.StatusDrifted, "drifted"},
		{drift.StatusMissing, "missing"},
		{drift.StatusOrphan, "orphan"},
	}
	for _, tc := range cases {
		if got := tc.status.String(); got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.status, got, tc.want)
		}
	}
}
