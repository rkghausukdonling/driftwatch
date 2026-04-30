// Package baseline manages drift scan baseline snapshots.
//
// A baseline snapshot captures the state of drift detection results at a
// specific point in time. Subsequent scans can be compared against the
// baseline to surface only newly introduced drift, reducing noise from
// pre-existing known issues.
//
// Typical usage:
//
//	// Save a baseline after a clean deployment:
//	_ = baseline.Save(".driftwatch.baseline.json", "aws", results)
//
//	// On the next scan, load and compare:
//	snap, _ := baseline.Load(".driftwatch.baseline.json")
//	newDrifts := baseline.Compare(snap, currentResults)
//
// Snapshots are stored as human-readable JSON files and can be committed
// to source control alongside your IaC definitions.
package baseline
