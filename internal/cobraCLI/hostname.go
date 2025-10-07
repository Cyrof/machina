package cobracli

import (
	"path/filepath"

	"github.com/Cyrof/machina/internal/run"
	"github.com/spf13/cobra"
)

var (
	newHostname string
	autoRestart bool
)

var hostnameCmd = &cobra.Command{
	Use:   "hostname",
	Short: "Change the computer's hostname",
	RunE: func(cmd *cobra.Command, args []string) error {
		if newHostname == "" {
			return cmd.Help()
		}
		ps := filepath.Join("scripts", "change-hostname.ps1")
		argsPS := []string{"-NewName", newHostname}
		if autoRestart {
			argsPS = append(argsPS, "-Restart")
		}
		return run.PS1(ps, argsPS...)
	},
}

func init() {
	hostnameCmd.Flags().StringVar(&newHostname, "name", "", "New hostname")
	hostnameCmd.Flags().BoolVar(&autoRestart, "restart", false, "Restart after changing the hostname")

	rootCmd.AddCommand(hostnameCmd)
}
