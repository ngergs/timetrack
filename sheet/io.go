package sheet

import (
	"encoding/json"
	"fmt"
	"github.com/ngergs/timetrack/v2/constants"
	"github.com/ngergs/timetrack/v2/sheet/states"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"
)

var filePattern, _ = regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2}\\.json")

func GetCurrentTimesheet(dir string, dailyWorkingMinutes int) (*Timesheet, error) {
	saved, err := getLastSavedTimesheet(dir)
	if err != nil {
		return nil, fmt.Errorf("could not load last saved timesheet: %w", err)
	}
	if saved == nil {
		log.Debug().Msg("No timesheet found, opening a new one")
		saved = &Timesheet{
			Slices: make([]Timeslice, 0),
		}
	} else {
		// validation of the timesheet guarantees that len(saved.Slices)>0
		if saved.Slices[0].Start.Format(constants.ReferenceFormat) != time.Now().Format(constants.ReferenceFormat) {
			saved = prepareNewFromOldSheet(saved, dailyWorkingMinutes, dir)
		}
	}
	return saved, nil
}

func prepareNewFromOldSheet(old *Timesheet, dailyWorkingMinutes int, saveFolder string) *Timesheet {
	log.Debug().Msg("Found old timesheet, opening a new one")
	new := &Timesheet{
		Balance: old.Balance + old.GetTodayBalance() - dailyWorkingMinutes,
		Slices:  make([]Timeslice, 0),
	}
	if old.GetState() == states.Open {
		log.Debug().Msg("Old sheet contains open session, continue it here")
		now := time.Now()
		new.Slices = append(old.Slices, Timeslice{Start: &now})
		new.Save(saveFolder)
		old.EndSession()
		old.Save(saveFolder)
	}
	return new
}

func getLastSavedTimesheet(dir string) (lastSaved *Timesheet, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// entries are sorted by filename, so start from the end
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if filePattern.MatchString(entry.Name()) &&
			entry.Type().IsRegular() {
			lastSaved, err = Read[Timesheet](path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			err = lastSaved.Validate()
			if err != nil {
				return nil, err
			}
			return lastSaved, nil
		}
	}
	return nil, nil
}

func Read[T any](path string) (*T, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result T
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func Write[T any](path string, content *T) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0755)
}
