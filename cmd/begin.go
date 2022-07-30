package cmd

import (
	"fmt"
	"github.com/ngergs/timetrack/v2/sheet"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "begin",
		Args:  cobra.NoArgs,
		Short: "Starts a new timetrack sessions",
		Long:  `This starts a new timetrack session or reports an error if a session is already open,`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := loadSheetAndDo(cmd, nil, true, func(timesheet *sheet.Timesheet) error {
				return timesheet.BeginSession()
			})
			if err != nil {
				return fmt.Errorf("failed to start a new time session: %w", err)
			}
			return nil
		},
	})
}
