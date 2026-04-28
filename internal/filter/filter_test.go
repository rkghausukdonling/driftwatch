package filter_test

import (
	"testing"

	"github.com/your-org/driftwatch/internal/drift"
	"github.com/your-org/driftwatch/internal/filter"
)

func makeResults() []drift.Result {
	return []drift.Result{
		{ResourceID: "bucket-1", ResourceType: "aws_s3_bucket", Status: drift.StatusMatch},
		{ResourceID: "bucket-2", ResourceType: "aws_s3_bucket", Status: drift.StatusDrifted},
		{ResourceID: "sg-abc", ResourceType: "aws_security_group", Status: drift.StatusMissing},
		{ResourceID: "sg-def", ResourceType: "aws_security_group", Status: drift.StatusDrifted},
	}
}

func TestApply_NoOptions(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{})
	if len(got) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(got))
	}
}

func TestApply_FilterByType(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{
		Types: []string{"aws_s3_bucket"},
	})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	for _, r := range got {
		if r.ResourceType != "aws_s3_bucket" {
			t.Errorf("unexpected type %q", r.ResourceType)
		}
	}
}

func TestApply_FilterByStatus(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{
		Statuses: []drift.Status{drift.StatusDrifted},
	})
	if len(got) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(got))
	}
	for _, r := range got {
		if r.Status != drift.StatusDrifted {
			t.Errorf("unexpected status %v", r.Status)
		}
	}
}

func TestApply_FilterByIDPrefix(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{
		IDPrefix: "sg-",
	})
	if len(got) != 2 {
		t.Fatalf("expected 2 sg results, got %d", len(got))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{
		Types:    []string{"aws_security_group"},
		Statuses: []drift.Status{drift.StatusDrifted},
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].ResourceID != "sg-def" {
		t.Errorf("expected sg-def, got %q", got[0].ResourceID)
	}
}

func TestApply_NoMatch(t *testing.T) {
	got := filter.Apply(makeResults(), filter.Options{
		Types: []string{"aws_lambda_function"},
	})
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}
