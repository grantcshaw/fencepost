package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/soyuz43/fencepost/internal/config"
	"github.com/soyuz43/fencepost/internal/keystore"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Manage arbitrary metadata fields for a service",
}

var metadataSetCmd = &cobra.Command{
	Use:   "set <service> <key> <value>",
	Short: "Set a metadata field",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		if err := s.SetMetadata(args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "metadata %q = %q set for %q\n", args[1], args[2], args[0])
		return nil
	},
}

var metadataGetCmd = &cobra.Command{
	Use:   "get <service> <key>",
	Short: "Get a metadata field",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		val, err := s.GetMetadata(args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, val)
		return nil
	},
}

var metadataListCmd = &cobra.Command{
	Use:   "list <service>",
	Short: "List all metadata fields for a service",
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
		m, err := s.AllMetadata(args[0])
		if err != nil {
			return err
		}
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(os.Stdout, "%s = %s\n", k, m[k])
		}
		return nil
	},
}

func init() {
	metadataCmd.AddCommand(metadataSetCmd, metadataGetCmd, metadataListCmd)
	rootCmd.AddCommand(metadataCmd)
}
