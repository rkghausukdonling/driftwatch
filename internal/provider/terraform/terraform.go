// Package terraform provides a provider that reads state from a Terraform state file.
package terraform

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/driftwatch/internal/provider"
)

const providerName = "terraform"

func init() {
	provider.Register(providerName, func(cfg map[string]string) (provider.Provider, error) {
		return newTerraformProvider(cfg)
	})
}

// terraformProvider reads resource state from a local Terraform state file.
type terraformProvider struct {
	statePath string
	resources map[string]map[string]string
}

// tfState mirrors the minimal structure of a terraform.tfstate file.
type tfState struct {
	Resources []tfResource `json:"resources"`
}

type tfResource struct {
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Instances []tfInstance `json:"instances"`
}

type tfInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

func newTerraformProvider(cfg map[string]string) (*terraformProvider, error) {
	path, ok := cfg["state_file"]
	if !ok || path == "" {
		return nil, fmt.Errorf("terraform provider: missing required config key 'state_file'")
	}

	p := &terraformProvider{statePath: path, resources: make(map[string]map[string]string)}
	if err := p.loadState(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *terraformProvider) loadState() error {
	data, err := os.ReadFile(p.statePath)
	if err != nil {
		return fmt.Errorf("terraform provider: reading state file: %w", err)
	}

	var state tfState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("terraform provider: parsing state file: %w", err)
	}

	for _, res := range state.Resources {
		for _, inst := range res.Instances {
			id, _ := inst.Attributes["id"].(string)
			if id == "" {
				continue
			}
			attrs := make(map[string]string)
			for k, v := range inst.Attributes {
				attrs[k] = fmt.Sprintf("%v", v)
			}
			attrs["resource_type"] = res.Type
			attrs["resource_name"] = res.Name
			p.resources[id] = attrs
		}
	}
	return nil
}

func (p *terraformProvider) Name() string { return providerName }

func (p *terraformProvider) FetchState(resourceID string) (map[string]string, error) {
	attrs, ok := p.resources[resourceID]
	if !ok {
		return nil, nil
	}
	return attrs, nil
}
