// Package remediate analyses drift detection results and produces
// actionable remediation suggestions for each non-OK resource.
//
// Usage:
//
//	results := detector.Detect(ctx, ids)
//	suggestions := remediate.Generate(results)
//	for _, s := range suggestions {
//		fmt.Println(s.Hint)
//	}
//
// Suggestions are only generated for resources whose status is
// [drift.StatusMissing] or [drift.StatusDrifted]; resources with
// [drift.StatusOK] are silently skipped.
package remediate
