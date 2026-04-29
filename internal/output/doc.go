// Package output provides utilities for selecting and constructing output
// destinations and formats for driftwatch scan results.
//
// Supported formats:
//
//	"text" – human-readable tabular output (default)
//	"json" – machine-readable JSON suitable for CI pipelines
//
// Usage:
//
//	fmt, err := output.ParseFormat(flagValue)
//	if err != nil {
//	    return err
//	}
//
//	w, err := output.WriterFor(outputPath) // "-" or "" for stdout
//	if err != nil {
//	    return err
//	}
//	defer w.Close()
//
//	// pass fmt and w to report.New()
package output
