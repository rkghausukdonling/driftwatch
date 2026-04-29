package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "driftwatch",
	Short: "Detect configuration drift between deployed infrastructure and IaC definitions",
	Long: `driftwatch compares your live infrastructure state against your
Infrastructure as Code definitions (Terraform, Pulumi, etc.) and reports
any configuration drift it detects.`,
	SilenceUsage: true,
}

// Execute runs the root command and exits with a non-zero status code on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// isVerbose returns true if verbose output has been enabled via the --verbose flag.
func isVerbose() bool {
	return verbose
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .driftwatch.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}
