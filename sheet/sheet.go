package sheet

import (
	"fmt"
	"github.com/ngergs/timetrack/v2/constants"
	"os"
	"text/tabwriter"
	"time"
)

func PrintStatus(sheet *Timesheet, date *time.Time, compactPrint bool) {
	statW := tabwriter.NewWriter(os.Stdout, 20, 20, 0, ' ', 0)
	fmt.Fprintf(statW, "Date:\t%s\n", date.Format(constants.DateOnlyFormat))
	fmt.Fprintf(statW, "Start date balance\t%dh%dmin\t\n", sheet.Balance/60, abs(sheet.Balance)%60)
	fmt.Fprintf(statW, "Worked today\t%dh%dmin\t\n", sheet.GetTodayBalance()/60, abs(sheet.GetTodayBalance())%60)
	fmt.Fprintf(statW, "Current session\t%s\t\n", sheet.GetState().String())
	statW.Flush()
	if !compactPrint {
		sliceW := tabwriter.NewWriter(os.Stdout, 28, 28, 0, ' ', 0)
		fmt.Fprintf(sliceW, "\nStart\tEnd\n")
		for _, slice := range sheet.Slices {
			fmt.Fprintf(sliceW, "%s\t", slice.Start.Format(constants.TimeFormat))
			if slice.End != nil {
				fmt.Fprintf(sliceW, "%s", slice.End.Format(constants.TimeFormat))
			}
			fmt.Fprint(sliceW, "\t\n")
			sliceW.Flush()
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
