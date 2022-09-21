package main

import (
	"fmt"
	"github.com/ngergs/timetrack/v2/internal/constants"
	"github.com/ngergs/timetrack/v2/sheet"
	"github.com/spf13/cobra"
	"time"
)

const compact = "compact"

func init() {
	command := &cobra.Command{
		Use:   "status <date, e.g. 2022-07-26>",
		Args:  cobra.MaximumNArgs(1),
		Short: "Prints the status of today's timesheet or from the specified date",
		Long: `Prints the status of today's timesheet if no date is specified.
If specified the timesheet from the specified date is printed.
Returns an error if no such timesheet exists.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			compact, err := cmd.Flags().GetBool(compact)
			if err != nil {
				return fmt.Errorf("failed to parse argument: %w", err)
			}
			var date *time.Time
			if len(args) == 1 {
				argDate, err := time.Parse(constants.DateOnlyFormat, args[0])
				if err != nil {
					return fmt.Errorf("invalid input date format: %w", err)
				}
				date = &argDate
			}
			// pass nil if no date is set to trigger the logic to load the last valid timesheet
			err = loadSheetAndDo(date, false, func(timesheet *sheet.Timesheet) error {
				if date == nil {
					now := time.Now()
					date = &now
				}
				sheet.PrintStatus(timesheet, date, compact)
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to print the time session status: %w", err)
			}
			return nil
		},
	}
	command.Flags().Bool(compact, false, "Whether to print the compact version of the sheet")
	rootCmd.AddCommand(command)
}
