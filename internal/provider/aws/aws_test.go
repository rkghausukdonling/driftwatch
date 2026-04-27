package aws

import (
	"context"
	"testing"

	"github.com/user/driftwatch/internal/provider"
)

func TestNewAWSProvider_MissingRegion(t *testing.T) {
	_, err := newAWSProvider(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing region, got nil")
	}
}

func TestNewAWSProvider_WithRegion(t *testing.T) {
	p, err := newAWSProvider(map[string]string{"region": "us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.region != "us-east-1" {
		t.Errorf("expected region us-east-1, got %q", p.region)
	}
	if p.profile != "" {
		t.Errorf("expected empty profile, got %q", p.profile)
	}
}

func TestNewAWSProvider_WithProfile(t *testing.T) {
	p, err := newAWSProvider(map[string]string{"region": "eu-west-1", "profile": "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.profile != "staging" {
		t.Errorf("expected profile 'staging', got %q", p.profile)
	}
}

func TestAWSProvider_Name(t *testing.T) {
	p, _ := newAWSProvider(map[string]string{"region": "us-west-2"})
	if p.Name() != ProviderName {
		t.Errorf("expected name %q, got %q", ProviderName, p.Name())
	}
}

func TestAWSProvider_FetchState_EmptyID(t *testing.T) {
	p, _ := newAWSProvider(map[string]string{"region": "us-east-1"})
	_, err := p.FetchState(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty resourceID, got nil")
	}
}

func TestAWSProvider_FetchState_NotImplemented(t *testing.T) {
	p, _ := newAWSProvider(map[string]string{"region": "us-east-1"})
	_, err := p.FetchState(context.Background(), "i-0abc123")
	if err == nil {
		t.Fatal("expected not-implemented error, got nil")
	}
}

func TestAWSProvider_RegisteredViaInit(t *testing.T) {
	available := provider.Available()
	for _, name := range available {
		if name == ProviderName {
			return
		}
	}
	t.Errorf("aws provider not found in Available(): %v", available)
}
