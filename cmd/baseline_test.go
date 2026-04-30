package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/baseline"
	"github.com/yourusername/driftwatch/internal/drift"
)

func TestBaselineCmd_HasSubcommands(t *testing.T) {
	subs := map[string]bool{}
	for _, c := range baselineCmd.Commands() {
		subs[c.Use] = true
	}
	for _, want := range []string{"save", "compare"} {
		if !subs[want] {
			t.Errorf("expected subcommand %q to be registered", want)
		}
	}
}

func TestBaselineCompareCmd_NoNewDrift(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	results := []drift.Result{
		{ID: "i-123", Type: "aws_instance", Status: drift.StatusMatch},
	}
	snap := baseline.Snapshot{
		CreatedAt: time.Now().UTC(),
		Provider:  "mock",
		Results:   results,
	}
	data, _ := json.MarshalIndent(snap, "", "  ")
	_ = os.WriteFile(path, data, 0o644)

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	diffs := baseline.Compare(loaded, results)
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(diffs))
	}
}

func TestBaselineCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "baseline" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'baseline' command to be registered on root")
	}
}

func TestBaselineSaveCmd_UsageOutput(t *testing.T) {
	var buf bytes.Buffer
	baselineSaveCmd.SetOut(&buf)
	if baselineSaveCmd.Use != "save" {
		t.Errorf("unexpected Use: %q", baselineSaveCmd.Use)
	}
}
