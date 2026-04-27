package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var reviewedCmd = &cobra.Command{
	Use:   "reviewed",
	Short: "Manage last-reviewed timestamps for services",
}

var reviewedSetCmd = &cobra.Command{
	Use:   "set <service>",
	Short: "Record the current time as the last review for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		if err := s.SetReviewedAt(args[0], time.Now()); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "reviewed timestamp set for %q\n", args[0])
		return nil
	},
}

var reviewedClearCmd = &cobra.Command{
	Use:   "clear <service>",
	Short: "Clear the review timestamp for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		return s.ClearReviewedAt(args[0])
	},
}

var reviewedNeverCmd = &cobra.Command{
	Use:   "never",
	Short: "List services that have never been reviewed",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		for _, name := range s.NeverReviewed() {
			fmt.Println(name)
		}
		return nil
	},
}

func init() {
	reviewedCmd.AddCommand(reviewedSetCmd)
	reviewedCmd.AddCommand(reviewedClearCmd)
	reviewedCmd.AddCommand(reviewedNeverCmd)
	AddCommand(reviewedCmd)
}
