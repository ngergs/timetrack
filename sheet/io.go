package sheet

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/ngergs/timetrack/v2/internal/constants"
	"github.com/ngergs/timetrack/v2/internal/io"
	"github.com/ngergs/timetrack/v2/sheet/states"
	"github.com/rs/zerolog/log"
)

var CurrentDayString string = time.Now().Format(constants.ReferenceFormat)
var filePattern, _ = regexp.Compile(`[0-9]{4}-[0-9]{2}-[0-9]{2}\.json`)

func GetLastTimesheet(lastSaved *os.File) (saved *Timesheet, err error) {
	if lastSaved == nil {
		log.Debug().Msg("No timesheet found, opening a new one")
		saved = &Timesheet{
			Slices: make([]Timeslice, 0),
		}
	} else {
		saved, err = io.Read[Timesheet](lastSaved)
		if err != nil {
			return nil, fmt.Errorf("could not load timesheet: %w", err)
		}
	}
	return saved, saved.Validate()
}

func PrepareNewFromOldSheet(old *Timesheet, dailyWorkingMinutes int) (*Timesheet, error) {
	log.Debug().Msg("Found old timesheet, opening a new one")
	current := &Timesheet{
		Balance: old.Balance + old.GetTodayBalance() - dailyWorkingMinutes,
		Slices:  make([]Timeslice, 0),
	}
	if old.GetState() == states.Open {
		log.Debug().Msg("Old sheet contains open session, copy it over")
		now := time.Now()
		current.Slices = append(old.Slices, Timeslice{Start: &now})
		err := old.EndSession()
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

// getLastSavedFile is an internal function used by GetLastSavedFile when no date is present.
func getLastSavedFile(folder string, write bool) (*os.File, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	// entries are sorted by filename, so start from the end
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if filePattern.MatchString(entry.Name()) &&
			entry.Type().IsRegular() {
			lastSaved, err := io.OpenFile(path.Join(folder, entry.Name()), write)
			if err != nil {
				return nil, err
			}
			return lastSaved, nil
		}
	}
	return nil, nil
}

// GetLastSavedFile returns the saved file for the given date if not nil or an error.
// If the date is nil the last saved file is determined and opened. If none is found a nil pointer without an error is returned.
func GetLastSavedFile(folder string, date *time.Time, write bool) (file *os.File, err error) {
	if date != nil {
		file, err = io.OpenFile(path.Join(folder, date.Format(constants.ReferenceFormat)), write)
		if err != nil {
			return nil, fmt.Errorf("failed to open last saved timesheet: %w", err)
		}
	} else {
		file, err = getLastSavedFile(folder, write)
		if err != nil {
			return nil, fmt.Errorf("failed to determine last saved timesheet: %w", err)
		}
	}
	return
}
