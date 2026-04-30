package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/driftwatch/internal/baseline"
	"github.com/yourusername/driftwatch/internal/config"
	"github.com/yourusername/driftwatch/internal/drift"
	"github.com/yourusername/driftwatch/internal/provider"
)

var baselineFile string

func init() {
	baselineCmd.Flags().StringVar(&baselineFile, "file", ".driftwatch.baseline.json", "path to baseline snapshot file")
	baselineCmd.AddCommand(baselineSaveCmd)
	baselineCmd.AddCommand(baselineCompareCmd)
	rootCmd.AddCommand(baselineCmd)
}

var baselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Manage drift scan baselines",
	Long:  "Save or compare drift scan results against a known-good baseline snapshot.",
}

var baselineSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current scan results as the baseline",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		p, err := provider.New(cfg.Provider, cfg.Options)
		if err != nil {
			return err
		}
		detector := drift.New(p)
		results, err := detector.Detect(cmd.Context(), cfg.ResourceIDs)
		if err != nil {
			return err
		}
		if err := baseline.Save(baselineFile, cfg.Provider, results); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Baseline saved to %s (%d resources)\n", baselineFile, len(results))
		return nil
	},
}

var baselineCompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare current scan against the saved baseline",
	RunE: func(cmd *cobra.Command, args []string) error {
		snap, err := baseline.Load(baselineFile)
		if err != nil {
			return err
		}
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		p, err := provider.New(cfg.Provider, cfg.Options)
		if err != nil {
			return err
		}
		detector := drift.New(p)
		results, err := detector.Detect(cmd.Context(), cfg.ResourceIDs)
		if err != nil {
			return err
		}
		diffs := baseline.Compare(snap, results)
		if len(diffs) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No new drift detected since baseline.")
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%d resource(s) have changed since baseline:\n", len(diffs))
		for _, r := range diffs {
			fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %s (%s)\n", r.Status, r.ID, r.Type)
		}
		os.Exit(1)
		return nil
	},
}
