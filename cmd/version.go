package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "snapshot"

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Args:  cobra.NoArgs,
		Short: "Prints the program version",
		Long: `Prints the program version.
The release notes can be found at: https://github.com/ngergs/timetrack/releases`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", version)
		},
	})
}
