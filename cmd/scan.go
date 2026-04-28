package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/report"

	// Register built-in providers.
	_ "github.com/driftwatch/internal/provider/aws"
	_ "github.com/driftwatch/internal/provider/mock"
	_ "github.com/driftwatch/internal/provider/terraform"
)

var (
	cfgFile    string
	outputFmt  string
	resourceIDs []string
)

// scanCmd represents the scan subcommand. It loads configuration, initialises
// the chosen provider, runs drift detection and writes a report to stdout.
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan infrastructure for configuration drift",
	Long: `Scan compares the live state of your infrastructure resources against
your IaC definitions and reports any drift that is detected.

Example:
  driftwatch scan --config driftwatch.yaml
  driftwatch scan --config driftwatch.yaml --output json
  driftwatch scan --config driftwatch.yaml --resource vpc-abc123 --resource sg-def456`,
	RunE: runScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&cfgFile, "config", "c", "driftwatch.yaml",
		"Path to the driftwatch configuration file")
	scanCmd.Flags().StringVarP(&outputFmt, "output", "o", "text",
		"Output format: text or json")
	scanCmd.Flags().StringArrayVarP(&resourceIDs, "resource", "r", nil,
		"Resource ID(s) to scan (overrides config resource_ids)")
}

func runScan(cmd *cobra.Command, args []string) error {
	// Load configuration from file.
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Command-line resource IDs take precedence over config.
	ids := cfg.ResourceIDs
	if len(resourceIDs) > 0 {
		ids = resourceIDs
	}

	if len(ids) == 0 {
		return fmt.Errorf("no resource IDs specified; set resource_ids in config or pass --resource flags")
	}

	// Initialise the provider specified in config.
	p, err := provider.New(cfg.Provider, cfg.ProviderConfig)
	if err != nil {
		return fmt.Errorf("initialising provider %q: %w", cfg.Provider, err)
	}

	// Run drift detection.
	detector := drift.New(p)
	results, err := detector.Detect(cmd.Context(), ids)
	if err != nil {
		return fmt.Errorf("running drift detection: %w", err)
	}

	// Write the report.
	r := report.New(results)

	switch outputFmt {
	case "json":
		if err := r.WriteJSON(os.Stdout); err != nil {
			return fmt.Errorf("writing JSON report: %w", err)
		}
	case "text":
		if err := r.WriteText(os.Stdout); err != nil {
			return fmt.Errorf("writing text report: %w", err)
		}
	default:
		return fmt.Errorf("unknown output format %q; supported formats: text, json", outputFmt)
	}

	// Exit with a non-zero status code when drift is detected so the command
	// can be used reliably in CI pipelines.
	for _, res := range results {
		if res.Status != "ok" {
			os.Exit(1)
		}
	}

	return nil
}
