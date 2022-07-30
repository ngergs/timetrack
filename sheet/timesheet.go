package sheet

import (
	"fmt"
	"github.com/ngergs/timetrack/v2/sheet/states"
	"math"
	"time"
)

type Timeslice struct {
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

type Timesheet struct {
	Balance int         `json:"balanceInMinutes"`
	Slices  []Timeslice `json:"slices"`
}

func newTimesheet() *Timesheet {
	return &Timesheet{
		Slices: make([]Timeslice, 0),
	}
}

func (sheet *Timesheet) BeginSession() error {
	if sheet.GetState() != states.Closed {
		return fmt.Errorf("current Timesheet has an open session, cannot start a new one")
	}
	now := time.Now()
	sheet.Slices = append(sheet.Slices, Timeslice{Start: &now})
	return nil
}

func (sheet *Timesheet) EndSession() error {
	if sheet.GetState() != states.Open {
		return fmt.Errorf("current Timesheet has no open session, therefore it cannot be closed")
	}
	now := time.Now()
	sheet.Slices[len(sheet.Slices)-1].End = &now
	return nil
}

func (sheet *Timesheet) GetState() states.State {
	if len(sheet.Slices) == 0 ||
		sheet.Slices[len(sheet.Slices)-1].End != nil {
		return states.Closed
	} else {
		return states.Open
	}
}

func (sheet *Timesheet) GetTodayBalance() int {
	var balance float64
	for _, entry := range sheet.Slices {
		if entry.End != nil {
			balance += entry.End.Sub(*entry.Start).Minutes()
		} else {
			balance += time.Now().Sub(*entry.Start).Minutes()
		}
	}
	return int(math.Floor(balance))
}

func (sheet *Timesheet) Validate() error {
	for i, entry := range sheet.Slices {
		if entry.Start == nil {
			return fmt.Errorf("time slice with empty start time found")
		}
		if entry.End == nil && i != len(sheet.Slices)-1 {
			return fmt.Errorf("time slice with empty end time found, but further time slices follow")
		}
	}
	return nil
}
