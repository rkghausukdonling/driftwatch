package provider_test

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/provider"
	_ "github.com/yourusername/driftwatch/internal/provider/mock" // side-effect: registers mock
)

func TestNew_UnknownProvider(t *testing.T) {
	_, err := provider.New("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for unknown provider, got nil")
	}
}

func TestNew_MockProvider(t *testing.T) {
	p, err := provider.New("mock", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "mock" {
		t.Errorf("expected name %q, got %q", "mock", p.Name())
	}
}

func TestMockProvider_FetchState_Found(t *testing.T) {
	p, _ := provider.New("mock", nil)

	resources, err := p.FetchState([]string{"instance-001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
	if resources[0].ID != "instance-001" {
		t.Errorf("unexpected resource ID: %s", resources[0].ID)
	}
}

func TestMockProvider_FetchState_NotFound(t *testing.T) {
	p, _ := provider.New("mock", nil)

	_, err := p.FetchState([]string{"does-not-exist"})
	if err == nil {
		t.Fatal("expected error for missing resource, got nil")
	}
}

func TestAvailable_ContainsMock(t *testing.T) {
	names := provider.Available()
	for _, n := range names {
		if n == "mock" {
			return
		}
	}
	t.Errorf("expected 'mock' in available providers, got %v", names)
}
