package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/user/fencepost/internal/config"
	"github.com/user/fencepost/internal/keystore"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search services by name, tag, or note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		store, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}

		query := args[0]
		results := store.Search(query)

		if len(results) == 0 {
			fmt.Printf("No services matched %q\n", query)
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tTAGS\tNOTE")
		for _, r := range results {
			tags := ""
			if len(r.Entry.Tags) > 0 {
				for i, tag := range r.Entry.Tags {
					if i > 0 {
						tags += ","
					}
					tags += tag
				}
			}
			note := r.Entry.Note
			if len(note) > 40 {
				note = note[:37] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.Service, tags, note)
		}
		return w.Flush()
	},
}

func init() {
	AddCommand(searchCmd)
}
