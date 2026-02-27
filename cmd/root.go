package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "honeycomb-cli",
	Short: "A CLI for the Honeycomb.io API",
	Long:  "A command-line interface for interacting with the Honeycomb.io observability platform.",
}

// Execute the root command and return an exit code.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func init() {
	rootCmd.PersistentFlags().String("api-key", "", "Honeycomb API key (or set HONEYCOMB_API_KEY)")
	rootCmd.PersistentFlags().String("api-url", "https://api.honeycomb.io", "Honeycomb API URL (or set HONEYCOMB_API_URL)")
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
