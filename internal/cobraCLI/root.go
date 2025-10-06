package cobracli

import (
	"fmt"
	"os"

	"github.com/Cyrof/machina/internal/elevate"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "machina",
	Short: "Windows system automation CLI (AD, hostname, network, etc)",
	Long:  "machina is a CLI tool to automate Windows system tasks using bat/ps1 scripts (AD, hostname, DNS, network, etc).",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !alreadyElevated && !elevate.IsAdmin() {
			fmt.Println("[*] Elevation required, relaunching as Administrator...")
			if err := elevate.RelaunchElevated(); err != nil {
				fmt.Fprintf(os.Stderr, "[x] Failed to elevate: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	},
}

var alreadyElevated bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
