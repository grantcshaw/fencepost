package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/clikd-inc/fencepost/internal/config"
	"github.com/clikd-inc/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var tierCmd = &cobra.Command{
	Use:   "tier",
	Short: "Manage service tier assignments",
}

var tierSetCmd = &cobra.Command{
	Use:   "set <service> <tier>",
	Short: "Set tier for a service (free, basic, pro, enterprise)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		if err := ks.SetTier(args[0], keystore.Tier(args[1])); err != nil {
			return err
		}
		fmt.Printf("tier set to %q for service %q\n", args[1], args[0])
		return nil
	},
}

var tierGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get tier for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		tier, err := ks.GetTier(args[0])
		if err != nil {
			return err
		}
		if tier == "" {
			fmt.Println("(none)")
		} else {
			fmt.Println(string(tier))
		}
		return nil
	},
}

var tierListCmd = &cobra.Command{
	Use:   "list <tier>",
	Short: "List services by tier",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		services := ks.ServicesByTier(keystore.Tier(args[0]))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tTIER")
		for _, s := range services {
			fmt.Fprintf(w, "%s\t%s\n", s, args[0])
		}
		w.Flush()
		return nil
	},
}

func init() {
	tierCmd.AddCommand(tierSetCmd, tierGetCmd, tierListCmd)
	AddCommand(tierCmd)
}
