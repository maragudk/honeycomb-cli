package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func newAuthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Verify your API key and show team/environment info",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := apiKey(cmd)
			if key == "" {
				return fmt.Errorf("API key is required (set HONEYCOMB_API_KEY or use --api-key)")
			}

			c := honeycomb.NewClient(key, honeycomb.WithBaseURL(apiURL(cmd)))
			auth, err := c.Auth(cmd.Context())
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Team:        %v\n", auth.Team.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Environment: %v\n", auth.Environment.Name)

			if len(auth.APIKeyAccess) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Permissions:")
				for perm, granted := range auth.APIKeyAccess {
					if granted {
						fmt.Fprintf(cmd.OutOrStdout(), "  %v\n", perm)
					}
				}
			}
			return nil
		},
	}
}
