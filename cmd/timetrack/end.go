package main

import (
	"fmt"
	"github.com/ngergs/timetrack/v2/sheet"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "end",
		Args:  cobra.NoArgs,
		Short: "Stops the current timetrack session",
		Long:  `This stops the currently running timetrack session or reports an error if none has been open.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := loadSheetAndDo(cmd, nil, true, func(timesheet *sheet.Timesheet) error {
				return timesheet.EndSession()
			})
			if err != nil {
				return fmt.Errorf("failed to end the time session: %w", err)
			}
			return nil
		},
	})
}
