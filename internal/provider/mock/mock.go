// Package mock provides a deterministic Provider implementation for use in
// tests and dry-run scenarios.
package mock

import (
	"fmt"

	"github.com/yourusername/driftwatch/internal/provider"
)

const Name = "mock"

func init() {
	provider.Register(Name, func(cfg map[string]string) (provider.Provider, error) {
		return &MockProvider{cfg: cfg, resources: defaultResources()}, nil
	})
}

// MockProvider satisfies provider.Provider with hard-coded data.
type MockProvider struct {
	cfg       map[string]string
	resources map[string]provider.Resource
}

func (m *MockProvider) Name() string { return Name }

func (m *MockProvider) FetchState(resourceIDs []string) ([]provider.Resource, error) {
	results := make([]provider.Resource, 0, len(resourceIDs))
	for _, id := range resourceIDs {
		r, ok := m.resources[id]
		if !ok {
			return nil, fmt.Errorf("mock: resource %q not found", id)
		}
		results = append(results, r)
	}
	return results, nil
}

// AddResource injects an extra resource, useful for table-driven tests.
func (m *MockProvider) AddResource(r provider.Resource) {
	m.resources[r.ID] = r
}

func defaultResources() map[string]provider.Resource {
	return map[string]provider.Resource{
		"instance-001": {
			ID:       "instance-001",
			Type:     "compute_instance",
			Provider: Name,
			Attributes: map[string]interface{}{
				"region":        "us-east-1",
				"instance_type": "t3.micro",
				"ami":           "ami-0abcdef1234567890",
			},
		},
	}
}
