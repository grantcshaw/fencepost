package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/fencepost/internal/keystore"
)

var mergeOverwrite bool

var mergeCmd = &cobra.Command{
	Use:   "merge <source-store>",
	Short: "Merge keys from another store file into the current store",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		dst := keystore.New(cfg.StorePath)
		src := keystore.New(args[0])

		results, err := dst.Merge(src, keystore.MergeOptions{Overwrite: mergeOverwrite})
		if err != nil {
			return fmt.Errorf("merge failed: %w", err)
		}

		if len(results) == 0 {
			fmt.Println("Nothing to merge.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tACTION")
		for _, r := range results {
			fmt.Fprintf(w, "%s\t%s\n", r.Service, r.Action)
		}
		w.Flush()
		return nil
	},
}

func init() {
	mergeCmd.Flags().BoolVar(&mergeOverwrite, "overwrite", false, "Overwrite existing keys in destination")
	AddCommand(mergeCmd)
}
