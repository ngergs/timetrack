package main

import (
	"github.com/ngergs/timetrack/modes"
	"github.com/ngergs/timetrack/states"
	"github.com/rs/zerolog/log"
	"time"
)

const referenceFormat string = "2006-01-02.json"

func main() {
	readConfig()
	sheet := getCurrentTimesheet(resolvedFolder)
	timesheetLogic(sheet)
}

func timesheetLogic(sheet *timesheet) {
	var err error
	switch mode {
	case modes.Status:
		log.Info().Msgf("Current session: %s", sheet.getState().String())
		log.Info().Msgf("Start day balance: %dh%dmin", sheet.Balance/60, abs(sheet.Balance)%60)
		log.Info().Msgf("Worked today: %dh%dmin", sheet.getTodayBalance()/60, abs(sheet.getTodayBalance())%60)
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
