package cmd

import (
	"fmt"
	"os"

	"github.com/fencepost/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	appConfig *config.Config
)

// rootCmd is the base command for the fencepost CLI.
var rootCmd = &cobra.Command{
	Use:   "fencepost",
	Short: "Manage and rotate API keys across services with audit logging",
	Long: `fencepost is a CLI tool that helps you manage API keys across
multiple services, rotate them on a schedule, and maintain a full
audit log of all key operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		appConfig, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "",
		"config file (default: $HOME/.fencepost/config.yaml)",
	)
}
