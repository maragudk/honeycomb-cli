package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags.
var version = "dev"

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})
}
