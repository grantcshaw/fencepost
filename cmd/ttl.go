package cmd

import (
	"fmt"
	"strconv"

	"github.com/smlrepo/fencepost/internal/config"
	"github.com/smlrepo/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Manage TTL (time-to-live) for service keys",
}

var ttlSetCmd = &cobra.Command{
	Use:   "set <service> <hours>",
	Short: "Set TTL in hours for a service key",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		hours, err := strconv.Atoi(args[1])
		if err != nil || hours < 0 {
			return fmt.Errorf("hours must be a non-negative integer")
		}
		s := keystore.New(cfg.StorePath)
		if err := s.SetTTL(args[0], hours); err != nil {
			return err
		}
		fmt.Printf("TTL for %q set to %d hours\n", args[0], hours)
		return nil
	},
}

var ttlGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get TTL for a service key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		hours, err := s.GetTTL(args[0])
		if err != nil {
			return err
		}
		if hours == 0 {
			fmt.Printf("%s: no TTL set\n", args[0])
		} else {
			fmt.Printf("%s: %d hours\n", args[0], hours)
		}
		return nil
	},
}

var ttlExpiredCmd = &cobra.Command{
	Use:   "expired",
	Short: "List services whose keys have exceeded their TTL",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		expired := s.ExpiredByTTL()
		if len(expired) == 0 {
			fmt.Println("No TTL-expired keys found.")
			return nil
		}
		for _, name := range expired {
			fmt.Println(name)
		}
		return nil
	},
}

func init() {
	ttlCmd.AddCommand(ttlSetCmd, ttlGetCmd, ttlExpiredCmd)
	AddCommand(ttlCmd)
}
