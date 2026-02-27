package cmd

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newTriggersCommand() *cobra.Command {
	triggersCmd := &cobra.Command{
		Use:   "triggers",
		Short: "Manage triggers (alerts)",
	}

	triggersCmd.PersistentFlags().String("dataset", "", "Dataset slug (required)")
	_ = triggersCmd.MarkPersistentFlagRequired("dataset")

	triggersCmd.AddCommand(newTriggersListCommand())
	triggersCmd.AddCommand(newTriggersGetCommand())

	return triggersCmd
}

func newTriggersListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all triggers",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")

			triggers, err := c.ListTriggers(cmd.Context(), dataset)
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(triggers)
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tTHRESHOLD\tDISABLED")
			for _, tr := range triggers {
				disabled := ""
				if tr.Disabled {
					disabled = "yes"
				}
				threshold := fmt.Sprintf("%v %v", tr.Threshold.Op, tr.Threshold.Value)
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", tr.ID, tr.Name, threshold, disabled)
			}
			return w.Flush()
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
	return cmd
}

func newTriggersGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a trigger by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")

			trigger, err := c.GetTrigger(cmd.Context(), dataset, args[0])
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(trigger)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "ID:          %v\n", trigger.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "Name:        %v\n", trigger.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Description: %v\n", trigger.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "Threshold:   %v %v\n", trigger.Threshold.Op, trigger.Threshold.Value)
			fmt.Fprintf(cmd.OutOrStdout(), "Frequency:   %vs\n", trigger.Frequency)
			disabled := "no"
			if trigger.Disabled {
				disabled = "yes"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Disabled:    %v\n", disabled)
			return nil
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
	return cmd
}
