package cmd

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newSLOsCommand() *cobra.Command {
	slosCmd := &cobra.Command{
		Use:   "slos",
		Short: "Manage SLOs",
	}

	slosCmd.PersistentFlags().String("dataset", "", "Dataset slug (required)")
	_ = slosCmd.MarkPersistentFlagRequired("dataset")

	slosCmd.AddCommand(newSLOsListCommand())
	slosCmd.AddCommand(newSLOsGetCommand())

	return slosCmd
}

func newSLOsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all SLOs",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")

			slos, err := c.ListSLOs(cmd.Context(), dataset)
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(slos)
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tTARGET\tPERIOD")
			for _, s := range slos {
				fmt.Fprintf(w, "%v\t%v\t%.2f%%\t%vd\n", s.ID, s.Name, s.TargetPercent(), s.TimePeriodDays)
			}
			return w.Flush()
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
	return cmd
}

func newSLOsGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an SLO by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")

			slo, err := c.GetSLO(cmd.Context(), dataset, args[0])
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(slo)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ID:          %v\n", slo.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Name:        %v\n", slo.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Description: %v\n", slo.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "Target:      %.2f%%\n", slo.TargetPercent())
			fmt.Fprintf(cmd.OutOrStdout(), "Period:      %v days\n", slo.TimePeriodDays)
			fmt.Fprintf(cmd.OutOrStdout(), "SLI:         %v\n", slo.SLI.Alias)
			return nil
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
	return cmd
}
