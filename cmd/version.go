package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Version   = "0.1.0"
	BuildDate = "2024-01-01"
	Commit    = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of driftwatch",
	Run: func(cmd *cobra.Command, args []string) {
		short, _ := cmd.Flags().GetBool("short")
		if short {
			fmt.Println(Version)
			return
		}
		fmt.Printf("driftwatch v%s\n", Version)
		fmt.Printf("  Build Date: %s\n", BuildDate)
		fmt.Printf("  Commit:     %s\n", Commit)
	},
}

func init() {
	versionCmd.Flags().Bool("short", false, "print only the version number")
	rootCmd.AddCommand(versionCmd)
}
