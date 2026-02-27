package cmd

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func newMarkersCommand() *cobra.Command {
	markersCmd := &cobra.Command{
		Use:   "markers",
		Short: "Manage markers (deploy annotations, etc.)",
	}

	markersCmd.PersistentFlags().String("dataset", "__all__", "Dataset slug (use __all__ for environment-wide)")

	markersCmd.AddCommand(newMarkersListCommand())
	markersCmd.AddCommand(newMarkersCreateCommand())
	markersCmd.AddCommand(newMarkersDeleteCommand())

	return markersCmd
}

func newMarkersListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List markers",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")

			markers, err := c.ListMarkers(cmd.Context(), dataset)
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(markers)
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tTYPE\tMESSAGE\tCREATED")
			for _, m := range markers {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", m.ID, m.Type, m.Message, m.CreatedAt)
			}
			return w.Flush()
		},
	}
	cmd.Flags().Bool("json", false, "Output as JSON")
	return cmd
}

func newMarkersCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a marker",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")

			markerType, _ := cmd.Flags().GetString("type")
			message, _ := cmd.Flags().GetString("message")
			url, _ := cmd.Flags().GetString("url")

			marker, err := c.CreateMarker(cmd.Context(), dataset, honeycomb.CreateMarkerRequest{
				Type:    markerType,
				Message: message,
				URL:     url,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created marker %v (type: %v)\n", marker.ID, marker.Type)
			return nil
		},
	}
	cmd.Flags().String("type", "", "Marker type (e.g. deploy)")
	cmd.Flags().String("message", "", "Marker message")
	cmd.Flags().String("url", "", "URL to associate with the marker")
	return cmd
}

func newMarkersDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a marker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")

			if err := c.DeleteMarker(cmd.Context(), dataset, args[0]); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Deleted marker %v\n", args[0])
			return nil
		},
	}
}
