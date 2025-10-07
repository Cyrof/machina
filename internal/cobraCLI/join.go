package cobracli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Cyrof/machina/internal/run"
	"github.com/spf13/cobra"
)

var (
	adDomain  string
	adUser    string
	adPass    string
	adPrompt  bool
	adRestart bool
	adDNS     string
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join this machine to an Active Directory domain",
	RunE: func(cmd *cobra.Command, args []string) error {
		if adDomain == "" {
			return fmt.Errorf("--domain is required")
		}
		if !adPrompt && (adUser == "" || adPass == "") {
			return fmt.Errorf("provide --prompt or both --user and --password")
		}

		ps := filepath.Join("scripts", "join-ad.ps1")
		argsPS := []string{
			"-Domain", adDomain,
		}
		if adPrompt {
			argsPS = append(argsPS, "-PromptForCredentials")
		} else {
			argsPS = append(argsPS, "-User", adUser, "-Password", adPass)
		}
		if adDNS != "" {
			parts := strings.Split(adDNS, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			argsPS = append(argsPS, "-DNSServer", strings.Join(parts, ","))
		}
		if adRestart {
			argsPS = append(argsPS, "-Restart")
		}
		return run.PS1(ps, argsPS...)
	},
}

func init() {
	joinCmd.Flags().StringVar(&adDomain, "domain", "", "Active Directory domain to join (e.g., 2d.com)")
	joinCmd.Flags().BoolVar(&adPrompt, "prompt", true, "Prompt for AD user credentials")
	joinCmd.Flags().StringVar(&adUser, "user", "", "AD user with permissions to join the domain (used if --prompt=false)")
	joinCmd.Flags().StringVar(&adPass, "password", "", "Password for the AD user (used if --prompt=false)")
	joinCmd.Flags().BoolVar(&adRestart, "restart", false, "Restart the machine after joining the domain")
	joinCmd.Flags().StringVar(&adDNS, "dns", "", "Comma-separated list of DNS servers to use (e.g., 192.168.100.60, 1.1.1.1)")

	rootCmd.AddCommand(joinCmd)
}
