package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/richbl/fencepost/internal/config"
	"github.com/richbl/fencepost/internal/keystore"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Manage per-service key change history",
}

var changelogListCmd = &cobra.Command{
	Use:   "list <service>",
	Short: "List changelog entries for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		entries, err := ks.GetChangelog(args[0])
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			fmt.Println("no changelog entries")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tEVENT\tDETAIL")
		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%s\t%s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Event, e.Detail)
		}
		return w.Flush()
	},
}

var changelogAddCmd = &cobra.Command{
	Use:   "add <service> <event> [detail]",
	Short: "Append a changelog entry for a service",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		detail := ""
		if len(args) == 3 {
			detail = args[2]
		}
		if err := ks.AppendChangelog(args[0], args[1], detail); err != nil {
			return err
		}
		fmt.Printf("changelog entry added for %q\n", args[0])
		return nil
	},
}

var changelogClearCmd = &cobra.Command{
	Use:   "clear <service>",
	Short: "Clear all changelog entries for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks, err := keystore.New(cfg.StorePath)
		if err != nil {
			return err
		}
		if err := ks.ClearChangelog(args[0]); err != nil {
			return err
		}
		fmt.Printf("changelog cleared for %q\n", args[0])
		return nil
	},
}

func init() {
	changelogCmd.AddCommand(changelogListCmd)
	changelogCmd.AddCommand(changelogAddCmd)
	changelogCmd.AddCommand(changelogClearCmd)
	AddCommand(changelogCmd)
}
