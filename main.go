package main

import (
	"fmt"
	"github.com/ngergs/timetrack/modes"
	"github.com/ngergs/timetrack/states"
	"github.com/rs/zerolog/log"
	"os"
	"text/tabwriter"
	"time"
)

const referenceFormat string = "2006-01-02.json"

func main() {
	readConfig()
	sheet := getCurrentTimesheet(resolvedFolder)
	timesheetLogic(sheet)
}

const timeFormat string = "2006-02-01 15:04:05 MST"

func timesheetLogic(sheet *timesheet) {
	var err error
	switch mode {
	case modes.Status:
		printStatus(sheet)
		return
	case modes.Start:
		err = sheet.StartSession()
	case modes.Stop:
		err = sheet.StopSession()
	}
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to update session")
	}
	err = sheet.Save()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to save updated timesheet")
	}
}

func getCurrentTimesheet(dir string) *timesheet {
	sheet, err := getLastSavedTimesheet(dir)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load last sheet")
	}
	if sheet == nil {
		log.Debug().Msg("No timesheet found, opening a new one")
		sheet = &timesheet{
			Slices: make([]timeslice, 0),
		}
	} else {
		// validation of the timesheet guarantees that len(sheet.Slices)>0
		if sheet.Slices[0].Start.Format(referenceFormat) != time.Now().Format(referenceFormat) {
			sheet = prepareNewFromOldSheet(sheet)
		}
	}
	return sheet
}

func prepareNewFromOldSheet(sheet *timesheet) *timesheet {
	log.Debug().Msg("Found old timesheet, opening a new one")
	oldSheet := sheet
	sheet = &timesheet{
		Balance: sheet.Balance + sheet.getTodayBalance() - *dailyWorkingMinutes,
		Slices:  make([]timeslice, 0),
	}
	if oldSheet.getState() == states.Open {
		log.Debug().Msg("Old sheet contains open session, continue it here")
		now := time.Now()
		sheet.Slices = append(sheet.Slices, timeslice{Start: &now})
		sheet.Save()
		oldSheet.StopSession()
		oldSheet.Save()
	}
	return sheet
}

func printStatus(sheet *timesheet) {
	statW := tabwriter.NewWriter(os.Stdout, 20, 20, 0, ' ', 0)
	fmt.Fprintf(statW, "Start day balance\t%dh%dmin\t\n", sheet.Balance/60, abs(sheet.Balance)%60)
	fmt.Fprintf(statW, "Worked today\t%dh%dmin\t\n", sheet.getTodayBalance()/60, abs(sheet.getTodayBalance())%60)
	fmt.Fprintf(statW, "Current session\t%s\t\n", sheet.getState().String())
	statW.Flush()
	if !*compactPrint {
		sliceW := tabwriter.NewWriter(os.Stdout, 28, 28, 0, ' ', 0)
		fmt.Fprintf(sliceW, "\nStart\tEnd\n")
		for _, slice := range sheet.Slices {
			fmt.Fprintf(sliceW, "%s\t", slice.Start.Format(timeFormat))
			if slice.End != nil {
				fmt.Fprintf(sliceW, "%s", slice.End.Format(timeFormat))
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
