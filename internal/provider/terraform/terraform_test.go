package terraform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempState(t *testing.T, state tfState) string {
	t.Helper()
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("marshal state: %v", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "terraform.tfstate")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write state file: %v", err)
	}
	return path
}

func TestNewTerraformProvider_MissingStateFile(t *testing.T) {
	_, err := newTerraformProvider(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing state_file config")
	}
}

func TestNewTerraformProvider_FileNotFound(t *testing.T) {
	_, err := newTerraformProvider(map[string]string{"state_file": "/nonexistent/path.tfstate"})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestTerraformProvider_Name(t *testing.T) {
	state := tfState{Resources: []tfResource{}}
	path := writeTempState(t, state)
	p, err := newTerraformProvider(map[string]string{"state_file": path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != providerName {
		t.Errorf("Name() = %q, want %q", p.Name(), providerName)
	}
}

func TestTerraformProvider_FetchState_Found(t *testing.T) {
	state := tfState{
		Resources: []tfResource{
			{
				Type: "aws_instance",
				Name: "web",
				Instances: []tfInstance{
					{Attributes: map[string]interface{}{"id": "i-abc123", "instance_type": "t3.micro"}},
				},
			},
		},
	}
	path := writeTempState(t, state)
	p, err := newTerraformProvider(map[string]string{"state_file": path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	attrs, err := p.FetchState("i-abc123")
	if err != nil {
		t.Fatalf("FetchState error: %v", err)
	}
	if attrs == nil {
		t.Fatal("expected attrs, got nil")
	}
	if attrs["instance_type"] != "t3.micro" {
		t.Errorf("instance_type = %q, want %q", attrs["instance_type"], "t3.micro")
	}
	if attrs["resource_type"] != "aws_instance" {
		t.Errorf("resource_type = %q, want %q", attrs["resource_type"], "aws_instance")
	}
}

func TestTerraformProvider_FetchState_NotFound(t *testing.T) {
	state := tfState{Resources: []tfResource{}}
	path := writeTempState(t, state)
	p, err := newTerraformProvider(map[string]string{"state_file": path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	attrs, err := p.FetchState("i-missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs != nil {
		t.Errorf("expected nil attrs for unknown resource, got %v", attrs)
	}
}
