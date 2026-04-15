package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"fencepost/internal/config"
	"fencepost/internal/keystore"
)

var tagCmd = &cobra.Command{
	Use:   "tag <service> <tag1,tag2,...>",
	Short: "Set tags on a stored API key",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		store, err := keystore.New(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}

		service := args[0]
		tags := splitTags(args[1])

		if err := store.SetTags(service, tags); err != nil {
			return fmt.Errorf("set tags for %q: %w", service, err)
		}

		fmt.Printf("Tags set for %q: %s\n", service, strings.Join(tags, ", "))
		return nil
	},
}

var filterTagCmd = &cobra.Command{
	Use:   "filter-tag <tag>",
	Short: "List services that have a specific tag",
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

		matches := store.FilterByTag(args[0])
		if len(matches) == 0 {
			fmt.Printf("No services tagged %q\n", args[0])
			return nil
		}
		for _, svc := range matches {
			fmt.Println(svc)
		}
		return nil
	},
}

func splitTags(raw string) []string {
	parts := strings.Split(raw, ",")
	var tags []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}
	return tags
}

func init() {
	AddCommand(tagCmd)
	AddCommand(filterTagCmd)
}
