package cobracli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cyrof/machina/internal/run"
	"github.com/spf13/cobra"
)

var (
	newHostname string
	autoRestart bool
	useRegistry bool
	skipConfirm bool
)

var hostnameCmd = &cobra.Command{
	Use:   "hostname",
	Short: "Change the computer's hostname",
	RunE: func(cmd *cobra.Command, args []string) error {
		if newHostname == "" {
			return cmd.Help()
		}

		if useRegistry && !skipConfirm {
			fmt.Println("[!] WARNING: You are forcing hostname change via registry.")
			fmt.Println("\t - AD will NOT be updated automatically.")
			fmt.Println("\t -Machine may stay domain-joined but name mismatch in AD")
			fmt.Println("Type YES to continue:")
			reader := bufio.NewReader(os.Stdin)
			inp, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToUpper(inp)) != "YES" {
				fmt.Println("Operation cancelled.")
				return nil
			}
		}

		ps := filepath.Join("scripts", "change-hostname.ps1")
		argsPS := []string{"-NewName", newHostname}
		if autoRestart {
			argsPS = append(argsPS, "-Restart")
		}
		if useRegistry {
			argsPS = append(argsPS, "-Registry")
		}
		return run.PS1(ps, argsPS...)
	},
}

func init() {
	hostnameCmd.Flags().StringVar(&newHostname, "name", "", "New hostname")
	hostnameCmd.MarkFlagRequired("name")

	hostnameCmd.Flags().BoolVar(&autoRestart, "restart", false, "Restart after changing the hostname")
	hostnameCmd.Flags().BoolVar(&useRegistry, "registry", false, "Force hostname change via registry (AD will NOT be updated)")
	hostnameCmd.Flags().BoolVar(&skipConfirm, "yes", false, "Skip confirmation prompt when using --registry")

	rootCmd.AddCommand(hostnameCmd)
}
