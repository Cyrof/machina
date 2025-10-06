package cobracli

import (
	"fmt"
	"path/filepath"
	"github.com/Cyrof/machina/internal/run"
	
	"github.com/spf13/cobra"
)

var workgroup string

var unjoinCmd = &cobra.Command{
	Use: "unjoin",
	Short: "Unjoin domain and join a workgroup",
	RunE: func(cmd *cobra.Command, args []string) error {
		if workgroup == "" {
			workgroup = "WORKGROUP"
		}
		script := filepath.Join("scripts", "unjoin_to_workgroup.bat")
		fmt.Printf("[*] Running unjoin -> workgroup (%s)\n", workgroup)
		return run.BAT(script, workgroup)
	},
}

func init() {
	unjoinCmd.Flags().StringVarP(&workgroup, "workgroup", "w", "WORKGROUP", "Workgroup name")
	rootCmd.AddCommand(unjoinCmd)
}