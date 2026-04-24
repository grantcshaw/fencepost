package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/fencepost/internal/config"
	"github.com/fencepost/internal/keystore"
)

var cipherCmd = &cobra.Command{
	Use:   "cipher",
	Short: "Manage encryption cipher for services",
}

var cipherSetCmd = &cobra.Command{
	Use:   "set <service> <cipher>",
	Short: "Set the cipher for a service (aes-256, aes-128, chacha20, none)",
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
		if err := s.SetCipher(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("cipher for %q set to %q\n", args[0], args[1])
		return nil
	},
}

var cipherGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get the cipher for a service",
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
		cipher, err := s.GetCipher(args[0])
		if err != nil {
			return err
		}
		fmt.Println(cipher)
		return nil
	},
}

var cipherListCmd = &cobra.Command{
	Use:   "list <cipher>",
	Short: "List services using a specific cipher",
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
		services, err := s.ServicesByCipher(args[0])
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tCIPHER")
		for _, svc := range services {
			fmt.Fprintf(w, "%s\t%s\n", svc, args[0])
		}
		return w.Flush()
	},
}

func init() {
	cipherCmd.AddCommand(cipherSetCmd, cipherGetCmd, cipherListCmd)
	AddCommand(cipherCmd)
}
