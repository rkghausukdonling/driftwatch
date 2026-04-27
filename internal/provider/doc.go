// Package provider defines the Provider interface and a central registry
// for infrastructure providers used by driftwatch.
//
// # Adding a new provider
//
// Create a sub-package (e.g. internal/provider/aws) and call
// provider.Register in its init() function:
//
//	func init() {
//		provider.Register("aws", func(cfg map[string]string) (provider.Provider, error) {
//			return newAWSProvider(cfg)
//		})
//	}
//
// Then blank-import the sub-package wherever providers are loaded
// (typically cmd/root.go or a providers.go bootstrap file).
//
// # Resource
//
// Resource is the canonical representation of a deployed infrastructure
// object returned by any provider. The Attributes map holds provider-
// specific fields and is compared against the IaC definition during
// drift detection.
package provider
