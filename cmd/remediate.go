package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/remediate"
	"github.com/spf13/cobra"
)

var remediateCmd = &cobra.Command{
	Use:   "remediate",
	Short: "Show remediation hints for drifted resources",
	Long: `Detect configuration drift and print actionable remediation
suggestions for every resource that is missing or drifted.`,
	RunE: runRemediate,
}

func init() {
	rootCmd.AddCommand(remediateCmd)
	remediateCmd.Flags().StringSliceP("ids", "i", nil, "Resource IDs to check (comma-separated)")
}

func runRemediate(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load("")
	if err != nil {
		cfg = config.DefaultConfig()
	}

	p, err := provider.New(cfg.Provider, cfg.ProviderOptions)
	if err != nil {
		return fmt.Errorf("initialising provider: %w", err)
	}

	ids, _ := cmd.Flags().GetStringSlice("ids")

	detector := drift.New(p)
	results, err := detector.Detect(cmd.Context(), ids)
	if err != nil {
		return fmt.Errorf("detecting drift: %w", err)
	}

	suggestions := remediate.Generate(results)
	if len(suggestions) == 0 {
		fmt.Fprintln(os.Stdout, "No remediation needed — all resources are in sync.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Found %d remediation suggestion(s):\n\n", len(suggestions))
	for i, s := range suggestions {
		fmt.Fprintf(os.Stdout, "%d. [%s] %s (%s)\n   %s\n\n",
			i+1, s.Status, s.ResourceID, s.ResourceType, s.Hint)
	}
	return nil
}
