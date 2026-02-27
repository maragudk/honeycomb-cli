package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

// NewRootCommand creates a new root command with all subcommands registered.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "honeycomb-cli",
		Short:         "A CLI for the Honeycomb.io API",
		Long:          "A command-line interface for interacting with the Honeycomb.io observability platform.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().String("api-key", "", "Honeycomb API key (or set HONEYCOMB_API_KEY)")
	root.PersistentFlags().String("api-url", "https://api.honeycomb.io", "Honeycomb API URL (or set HONEYCOMB_API_URL)")

	root.AddCommand(newVersionCommand())
	root.AddCommand(newAuthCommand())
	root.AddCommand(newDatasetsCommand())
	root.AddCommand(newMarkersCommand())
	root.AddCommand(newColumnsCommand())
	root.AddCommand(newQueryCommand())

	return root
}

// Execute the root command and return an exit code.
func Execute() int {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// apiKey returns the API key from the flag or environment variable.
func apiKey(cmd *cobra.Command) string {
	key, _ := cmd.Flags().GetString("api-key")
	if key != "" {
		return key
	}
	return os.Getenv("HONEYCOMB_API_KEY")
}

// apiURL returns the API URL from the flag or environment variable.
func apiURL(cmd *cobra.Command) string {
	url, _ := cmd.Flags().GetString("api-url")
	if url != "https://api.honeycomb.io" {
		return url
	}
	if envURL := os.Getenv("HONEYCOMB_API_URL"); envURL != "" {
		return envURL
	}
	return url
}

// newClient creates a new Honeycomb API client from the command's flags.
func newClient(cmd *cobra.Command) *honeycomb.Client {
	return honeycomb.NewClient(apiKey(cmd), honeycomb.WithBaseURL(apiURL(cmd)))
}
