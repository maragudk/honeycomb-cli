package cmd

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func newDatasetsCommand() *cobra.Command {
	datasetsCmd := &cobra.Command{
		Use:   "datasets",
		Short: "Manage datasets",
	}

	datasetsCmd.AddCommand(newDatasetsListCommand())
	datasetsCmd.AddCommand(newDatasetsGetCommand())
	datasetsCmd.AddCommand(newDatasetsCreateCommand())
	datasetsCmd.AddCommand(newDatasetsDeleteCommand())

	return datasetsCmd
}

func newDatasetsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all datasets",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)

			datasets, err := c.ListDatasets(cmd.Context())
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(datasets)
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSLUG\tDESCRIPTION\tLAST WRITTEN")
			for _, d := range datasets {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", d.Name, d.Slug, d.Description, d.LastWrittenAt)
			}
			return w.Flush()
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
	return cmd
}

func newDatasetsGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <slug>",
		Short: "Get a dataset by slug",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)

			dataset, err := c.GetDataset(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(dataset)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Name:         %v\n", dataset.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Slug:         %v\n", dataset.Slug)
			fmt.Fprintf(cmd.OutOrStdout(), "Description:  %v\n", dataset.Description)
			fmt.Fprintf(cmd.OutOrStdout(), "Last written: %v\n", dataset.LastWrittenAt)
			return nil
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
	return cmd
}

func newDatasetsCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a dataset",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)

			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")

			dataset, err := c.CreateDataset(cmd.Context(), honeycomb.CreateDatasetRequest{
				Name:        name,
				Description: description,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created dataset %q (%v)\n", dataset.Name, dataset.Slug)
			return nil
		},
	}
	cmd.Flags().String("name", "", "Dataset name")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().String("description", "", "Dataset description")
	return cmd
}

func newDatasetsDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <slug>",
		Short: "Delete a dataset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)

			if err := c.DeleteDataset(cmd.Context(), args[0]); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Deleted dataset %q\n", args[0])
			return nil
		},
	}
}
