package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yourusername/fencepost/internal/config"
	"github.com/yourusername/fencepost/internal/keystore"
)

var costCmd = &cobra.Command{
	Use:   "cost",
	Short: "Manage monthly cost estimates for API keys",
}

var costSetCmd = &cobra.Command{
	Use:   "set <service> <cents>",
	Short: "Set the monthly cost estimate in USD cents",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		cents, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid cents value: %w", err)
		}
		return ks.SetCost(args[0], cents)
	},
}

var costGetCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Get the monthly cost estimate in USD cents",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		cost, err := ks.GetCost(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s: %d cents/month ($%.2f)\n", args[0], cost, float64(cost)/100)
		return nil
	},
}

var costAboveCmd = &cobra.Command{
	Use:   "above <cents>",
	Short: "List services with cost above threshold",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return err
		}
		ks := keystore.New(cfg.StorePath)
		threshold, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid threshold: %w", err)
		}
		results := ks.ServicesByCostAbove(threshold)
		if len(results) == 0 {
			fmt.Fprintln(os.Stdout, "no services above threshold")
			return nil
		}
		for _, svc := range results {
			cost, _ := ks.GetCost(svc)
			fmt.Fprintf(os.Stdout, "%-30s %d cents ($%.2f)\n", svc, cost, float64(cost)/100)
		}
		return nil
	},
}

func init() {
	costCmd.AddCommand(costSetCmd, costGetCmd, costAboveCmd)
	AddCommand(costCmd)
}
