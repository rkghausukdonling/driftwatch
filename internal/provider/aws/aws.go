// Package aws provides an AWS infrastructure provider for driftwatch.
// It fetches live resource state from AWS using the AWS SDK.
package aws

import (
	"context"
	"fmt"

	"github.com/user/driftwatch/internal/provider"
)

const ProviderName = "aws"

// Resource represents a simplified AWS resource state.
type Resource struct {
	ID     string
	Type   string
	Region string
	Tags   map[string]string
	Raw    map[string]interface{}
}

// awsProvider implements provider.Provider for AWS.
type awsProvider struct {
	region  string
	profile string
}

// Config holds AWS-specific provider configuration.
type Config struct {
	Region  string `yaml:"region"`
	Profile string `yaml:"profile"`
}

func init() {
	provider.Register(ProviderName, func(opts map[string]string) (provider.Provider, error) {
		return newAWSProvider(opts)
	})
}

func newAWSProvider(opts map[string]string) (*awsProvider, error) {
	region, ok := opts["region"]
	if !ok || region == "" {
		return nil, fmt.Errorf("aws provider: missing required option 'region'")
	}
	profile := opts["profile"] // optional
	return &awsProvider{region: region, profile: profile}, nil
}

// FetchState retrieves the current state of an AWS resource by its ID.
// In production this would call the appropriate AWS service API.
func (p *awsProvider) FetchState(ctx context.Context, resourceID string) (map[string]interface{}, error) {
	if resourceID == "" {
		return nil, fmt.Errorf("aws provider: resourceID must not be empty")
	}
	// Placeholder: real implementation would dispatch based on resource type
	// (e.g., EC2 DescribeInstances, S3 GetBucketTagging, etc.)
	return nil, fmt.Errorf("aws provider: FetchState not yet implemented for resource %q", resourceID)
}

// Name returns the canonical provider name.
func (p *awsProvider) Name() string {
	return ProviderName
}
