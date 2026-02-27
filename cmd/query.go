package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func newQueryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Run a query against a dataset",
		Long: `Run a query against a dataset and display results.

Examples:
  # Count all events in the last 2 hours
  honeycomb-cli query --dataset requests --calculation COUNT

  # Average duration broken down by status code
  honeycomb-cli query --dataset requests --calculation "AVG:duration_ms" --breakdown status_code

  # P99 duration with a filter
  honeycomb-cli query --dataset requests --calculation "P99:duration_ms" --filter "status_code = 200"

  # Multiple calculations
  honeycomb-cli query --dataset requests --calculation COUNT --calculation "AVG:duration_ms"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := newClient(cmd)
			dataset, _ := cmd.Flags().GetString("dataset")
			calcs, _ := cmd.Flags().GetStringSlice("calculation")
			breakdowns, _ := cmd.Flags().GetStringSlice("breakdown")
			filters, _ := cmd.Flags().GetStringSlice("filter")
			timeRange, _ := cmd.Flags().GetInt("time-range")
			limit, _ := cmd.Flags().GetInt("limit")

			spec := honeycomb.QuerySpec{
				TimeRange: timeRange,
			}

			if limit > 0 {
				spec.Limit = limit
			}

			for _, calc := range calcs {
				c, err := ParseCalculation(calc)
				if err != nil {
					return err
				}
				spec.Calculations = append(spec.Calculations, c)
			}

			if len(spec.Calculations) == 0 {
				spec.Calculations = []honeycomb.Calculation{{Op: "COUNT"}}
			}

			spec.Breakdowns = breakdowns

			for _, f := range filters {
				filter, err := ParseFilter(f)
				if err != nil {
					return err
				}
				spec.Filters = append(spec.Filters, filter)
			}

			result, err := c.RunQuery(cmd.Context(), dataset, spec)
			if err != nil {
				return err
			}

			asJSON, _ := cmd.Flags().GetBool("json")
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result.Data.Results)
			}

			if len(result.Data.Results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No results.")
				return nil
			}

			return printQueryResults(cmd, result.Data.Results)
		},
	}

	cmd.Flags().String("dataset", "", "Dataset slug (required)")
	_ = cmd.MarkFlagRequired("dataset")
	cmd.Flags().StringSlice("calculation", nil, "Calculation (e.g. COUNT, AVG:column, P99:column)")
	cmd.Flags().StringSlice("breakdown", nil, "Breakdown column")
	cmd.Flags().StringSlice("filter", nil, "Filter (e.g. \"status_code = 200\")")
	cmd.Flags().Int("time-range", 7200, "Time range in seconds (default 2 hours)")
	cmd.Flags().Int("limit", 0, "Maximum number of results")
	cmd.Flags().Bool("json", false, "Output as JSON")

	return cmd
}

// ParseCalculation from a string like "COUNT" or "AVG:duration_ms".
func ParseCalculation(s string) (honeycomb.Calculation, error) {
	parts := strings.SplitN(s, ":", 2)
	op := strings.ToUpper(parts[0])

	switch op {
	case "COUNT":
		return honeycomb.Calculation{Op: op}, nil
	case "SUM", "AVG", "COUNT_DISTINCT", "MAX", "MIN", "P001", "P01", "P05",
		"P10", "P25", "P50", "P75", "P90", "P95", "P99", "P999",
		"RATE_AVG", "RATE_SUM", "RATE_MAX":
		if len(parts) < 2 || parts[1] == "" {
			return honeycomb.Calculation{}, fmt.Errorf("calculation %v requires a column (e.g. %v:column_name)", op, op)
		}
		return honeycomb.Calculation{Op: op, Column: parts[1]}, nil
	default:
		return honeycomb.Calculation{}, fmt.Errorf("unknown calculation operator: %v", op)
	}
}

// ParseFilter from a string like "status_code = 200".
func ParseFilter(s string) (honeycomb.Filter, error) {
	operators := []string{"!=", ">=", "<=", "=", ">", "<", "contains", "does-not-contain",
		"starts-with", "does-not-start-with", "exists", "does-not-exist",
		"in", "not-in"}

	for _, op := range operators {
		parts := strings.SplitN(s, " "+op+" ", 2)
		if len(parts) == 2 {
			return honeycomb.Filter{
				Column: strings.TrimSpace(parts[0]),
				Op:     op,
				Value:  strings.TrimSpace(parts[1]),
			}, nil
		}

		// Handle unary operators (exists, does-not-exist)
		trimmed := strings.TrimSpace(s)
		if strings.HasSuffix(trimmed, " "+op) {
			return honeycomb.Filter{
				Column: strings.TrimSuffix(trimmed, " "+op),
				Op:     op,
			}, nil
		}
	}

	return honeycomb.Filter{}, fmt.Errorf("could not parse filter: %q (expected format: \"column op value\")", s)
}

func printQueryResults(cmd *cobra.Command, results []map[string]any) error {
	if len(results) == 0 {
		return nil
	}

	// Collect column headers from first result
	var headers []string
	for key := range results[0] {
		headers = append(headers, key)
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, row := range results {
		var vals []string
		for _, h := range headers {
			vals = append(vals, fmt.Sprint(row[h]))
		}
		fmt.Fprintln(w, strings.Join(vals, "\t"))
	}
	return w.Flush()
}
