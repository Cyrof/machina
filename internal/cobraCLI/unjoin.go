package cobracli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Cyrof/machina/internal/run"

	"github.com/spf13/cobra"
)

var (
	workgroup string
	restart   bool
)

var unjoinCmd = &cobra.Command{
	Use:   "unjoin",
	Short: "Unjoin domain and join a workgroup",
	RunE: func(cmd *cobra.Command, args []string) error {
		if workgroup == "" {
			workgroup = "WORKGROUP"
		}

		fmt.Println("[!] WARNING: This will unjoin the machine from the domain.")
		fmt.Println("[!] Make sure you can log in with a local Administrator account, otherwise you may loss access.")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Type YES to continue: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))

		if input != "YES" {
			fmt.Println("Operation cancelled.")
			return nil
		}

		psArgs := []string{"-Workgroup", workgroup}
		if restart {
			psArgs = append(psArgs, "-Restart")
		}
		return run.PS1Embedded("unjoin_to_workgroup.ps1", psArgs...)
	},
}

func init() {
	unjoinCmd.Flags().StringVarP(&workgroup, "workgroup", "w", "WORKGROUP", "Workgroup name")
	unjoinCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Restart the machine after unjoining")
	rootCmd.AddCommand(unjoinCmd)
}
