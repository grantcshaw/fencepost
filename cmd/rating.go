package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/fencepost/internal/config"
	"github.com/nicholasgasior/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var ratingCmd = &cobra.Command{
	Use:   "rating",
	Short: "Manage service ratings (critical, high, medium, low)",
}

var ratingSetCmd = &cobra.Command{
	Use:   "set <service> <rating>",
	Short: "Set rating for a service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		if err := s.SetRating(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Rating for %q set to %s\n", args[0], args[1])
		return nil
	},
}

var ratingGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get rating for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		r, err := s.GetRating(args[0])
		if err != nil {
			return err
		}
		fmt.Println(r)
		return nil
	},
}

var ratingListCmd = &cobra.Command{
	Use:   "list <rating>",
	Short: "List services by rating",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		s := keystore.New(cfg.StorePath)
		svcs, err := s.ServicesByRating(args[0])
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tRATING")
		for _, svc := range svcs {
			fmt.Fprintf(w, "%s\t%s\n", svc, args[0])
		}
		return w.Flush()
	},
}

func init() {
	ratingCmd.AddCommand(ratingSetCmd, ratingGetCmd, ratingListCmd)
	AddCommand(ratingCmd)
}
