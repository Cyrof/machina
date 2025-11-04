package cobracli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Cyrof/machina/internal/run"
	"github.com/spf13/cobra"
)

var (
	cleanCfgPath string
	cleanForce   bool
	cleanDryRun  bool
	cleanVerbose bool
)

type CleanupConfig struct {
	LogDirs    []string `json:"LogDirs"`
	TempFiles  []string `json:"TempFiles"`
	RegKeys    []string `json:"RegKeys"`
	Services   []string `json:"Services"`
	ExtraPaths []string `json:"ExtraPaths"`
}

func readCleanupConfig(path string) (*CleanupConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg CleanupConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}
	return &cfg, nil
}

func summariseTargets(cfg *CleanupConfig) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Targets summary:\n")
	fmt.Fprintf(&b, "\tLogDirs: %d\n", len(cfg.LogDirs))
	fmt.Fprintf(&b, "\tTempFiles: %d", len(cfg.TempFiles))
	fmt.Fprintf(&b, "\tRegKeys: %d\n", len(cfg.RegKeys))
	fmt.Fprintf(&b, "\tServices: %d\n", len(cfg.Services))
	fmt.Fprintf(&b, "\tExtraPaths: %d\n", len(cfg.ExtraPaths))
	return b.String()
}

func confirmPassword(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	ans, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	ans = strings.TrimSpace(strings.ToLower(ans))
	return ans == "y" || ans == "yes", nil
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleanup Machina artifacts (logs, registry, temp files) based on a config file",
	Example: ` machina clean -f C:\machina\cleanup.json
 machina clean -f cleanup.json --dry-run
 machina clean -f cleanup.json --force --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cleanCfgPath == "" {
			return fmt.Errorf("--file is required")
		}
		if _, err := os.Stat(cleanCfgPath); err != nil {
			return fmt.Errorf("config file not found or unreadable: %s", cleanCfgPath)
		}

		cfg, err := readCleanupConfig(cleanCfgPath)
		if err != nil {
			return err
		}

		fmt.Println(summariseTargets(cfg))
		if cleanVerbose {
			if len(cfg.LogDirs) > 0 {
				fmt.Println("LogDirs:")
				for _, v := range cfg.LogDirs {
					fmt.Printf("\t- %s\n", v)
				}
			}
			if len(cfg.TempFiles) > 0 {
				fmt.Println("TempFiles:")
				for _, v := range cfg.TempFiles {
					fmt.Printf("\t- %s\n", v)
				}
			}
			if len(cfg.RegKeys) > 0 {
				fmt.Println("RegKeys:")
				for _, v := range cfg.RegKeys {
					fmt.Printf("\t- %s\n", v)
				}
			}
			if len(cfg.Services) > 0 {
				fmt.Println("Services:")
				for _, v := range cfg.Services {
					fmt.Printf("\t- %s\n", v)
				}
			}
			if len(cfg.ExtraPaths) > 0 {
				fmt.Println("ExtraPaths:")
				for _, v := range cfg.ExtraPaths {
					fmt.Printf("\t- %s\n", v)
				}
			}
			fmt.Println()
		}

		// confirmation handlers
		if !cleanForce && !cleanDryRun {
			ok, err := confirmPassword("Proceed with cleanup? (y/N): ")
			if err != nil {
				return fmt.Errorf("configuration error: %w", err)
			}
			if !ok {
				fmt.Println("Cleanup aborted by user.")
				return nil
			}
		}

		argsPS := []string{
			"-ConfigPath", cleanCfgPath,
		}
		if cleanDryRun {
			argsPS = append(argsPS, "-DryRun")
		}
		if cleanVerbose {
			argsPS = append(argsPS, "-VerboseMode")
		}

		return run.PS1Embedded("cleanup.ps1", argsPS...)
	},
}

func init() {
	cleanCmd.Flags().StringVarP(&cleanCfgPath, "file", "f", "", "Path to cleanup config JSON (required)")
	cleanCmd.Flags().BoolVar(&cleanForce, "force", false, "Skip confirmation prompts")
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "Show what would be removed without deleting")
	cleanCmd.Flags().BoolVar(&cleanVerbose, "verbose", false, "Verbose output")
	_ = cleanCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(cleanCmd)

}
