package cmd

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newColumnsCommand() *cobra.Command {
	columnsCmd := &cobra.Command{
		Use:   "columns",
		Short: "Manage columns in a dataset",
	}

	columnsCmd.PersistentFlags().String("dataset", "", "Dataset slug (required)")
	_ = columnsCmd.MarkPersistentFlagRequired("dataset")

	columnsCmd.AddCommand(newColumnsListCommand())

	return columnsCmd
}

func newColumnsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List columns in a dataset",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")

			columns, err := c.ListColumns(cmd.Context(), dataset)
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(columns)
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "KEY NAME\tTYPE\tDESCRIPTION\tHIDDEN")
			for _, col := range columns {
				hidden := ""
				if col.Hidden {
					hidden = "yes"
				}
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", col.KeyName, col.Type, col.Description, hidden)
			}
			return w.Flush()
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
	return cmd
}
