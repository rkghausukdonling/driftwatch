package provider

import (
	"fmt"
	"sort"
)

// Provider defines the interface for infrastructure providers
// that can fetch the current state of deployed resources.
type Provider interface {
	// Name returns the provider identifier (e.g. "aws", "gcp").
	Name() string
	// FetchState retrieves the current state of resources matching the given IDs.
	FetchState(resourceIDs []string) ([]Resource, error)
}

// Resource represents a deployed infrastructure resource.
type Resource struct {
	ID         string
	Type       string
	Provider   string
	Attributes map[string]interface{}
}

// Registry holds registered provider factories.
var registry = map[string]func(cfg map[string]string) (Provider, error){}

// Register adds a provider factory under the given name.
func Register(name string, factory func(cfg map[string]string) (Provider, error)) {
	registry[name] = factory
}

// New returns a Provider for the given name, initialised with cfg.
func New(name string, cfg map[string]string) (Provider, error) {
	factory, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider %q: did you forget to import it?", name)
	}
	return factory(cfg)
}

// Available returns the names of all registered providers in sorted order.
func Available() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
