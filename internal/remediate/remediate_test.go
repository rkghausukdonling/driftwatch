package remediate_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/remediate"
)

func makeResult(id, typ string, status drift.Status, desired, actual map[string]string) drift.Result {
	return drift.Result{
		ID:      id,
		Type:    typ,
		Status:  status,
		Desired: desired,
		Actual:  actual,
	}
}

func TestGenerate_SkipsOK(t *testing.T) {
	results := []drift.Result{
		makeResult("res-1", "aws_instance", drift.StatusOK, nil, nil),
	}
	got := remediate.Generate(results)
	if len(got) != 0 {
		t.Fatalf("expected 0 suggestions for OK result, got %d", len(got))
	}
}

func TestGenerate_MissingResource(t *testing.T) {
	results := []drift.Result{
		makeResult("res-2", "aws_s3_bucket", drift.StatusMissing, nil, nil),
	}
	got := remediate.Generate(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	s := got[0]
	if s.ResourceID != "res-2" {
		t.Errorf("unexpected ResourceID: %s", s.ResourceID)
	}
	if !strings.Contains(s.Hint, "not found in the provider") {
		t.Errorf("hint missing expected text: %s", s.Hint)
	}
}

func TestGenerate_DriftedResource_WithFields(t *testing.T) {
	desired := map[string]string{"instance_type": "t3.micro", "ami": "ami-123"}
	actual := map[string]string{"instance_type": "t3.large", "ami": "ami-123"}
	results := []drift.Result{
		makeResult("res-3", "aws_instance", drift.StatusDrifted, desired, actual),
	}
	got := remediate.Generate(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Hint, "instance_type") {
		t.Errorf("hint should mention drifted field, got: %s", got[0].Hint)
	}
}

func TestGenerate_DriftedResource_NoFields(t *testing.T) {
	results := []drift.Result{
		makeResult("res-4", "aws_instance", drift.StatusDrifted, nil, nil),
	}
	got := remediate.Generate(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if !strings.Contains(got[0].Hint, "drifted from its IaC definition") {
		t.Errorf("unexpected hint: %s", got[0].Hint)
	}
}

func TestGenerate_Mixed(t *testing.T) {
	results := []drift.Result{
		makeResult("ok-1", "aws_instance", drift.StatusOK, nil, nil),
		makeResult("miss-1", "aws_s3_bucket", drift.StatusMissing, nil, nil),
		makeResult("drift-1", "aws_instance", drift.StatusDrifted, map[string]string{"x": "a"}, map[string]string{"x": "b"}),
	}
	got := remediate.Generate(results)
	if len(got) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(got))
	}
}
