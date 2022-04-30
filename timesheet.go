package main

import (
	"fmt"
	"github.com/ngergs/timetrack/states"
	"math"
	"path"
	"time"
)

type timeslice struct {
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

type timesheet struct {
	Balance int         `json:"balanceInMinutes"`
	Slices  []timeslice `json:"slices"`
}

func newTimesheet() *timesheet {
	return &timesheet{
		Slices: make([]timeslice, 0),
	}
}

func (sheet *timesheet) StartSession() error {
	if sheet.getState() != states.Closed {
		return fmt.Errorf("current timesheet has an open session, cannot start a new one")
	}
	now := time.Now()
	sheet.Slices = append(sheet.Slices, timeslice{Start: &now})
	return nil
}

func (sheet *timesheet) StopSession() error {
	if sheet.getState() != states.Open {
		return fmt.Errorf("current timesheet has no open session, therefore it cannot be closed")
	}
	now := time.Now()
	sheet.Slices[len(sheet.Slices)-1].End = &now
	return nil
}

func (sheet *timesheet) getState() states.State {
	if len(sheet.Slices) == 0 ||
		sheet.Slices[len(sheet.Slices)-1].End != nil {
		return states.Closed
	} else {
		return states.Open
	}
}

func (sheet *timesheet) getTodayBalance() int {
	balance := 0
	for _, entry := range sheet.Slices {
		if entry.End != nil {
			balance += int(math.Floor(entry.End.Sub(*entry.Start).Minutes()))
		} else {
			balance += int(math.Floor(time.Now().Sub(*entry.Start).Minutes()))
		}
	}
	return balance
}

func (sheet *timesheet) Save() error {
	savename := sheet.Slices[0].Start.Format(referenceFormat)
	return write(path.Join(resolvedFolder, savename), sheet)
}

func (sheet *timesheet) Validate() error {
	if len(sheet.Slices) == 0 {
		return fmt.Errorf("empty sheets should not occur during normal operation")
	}
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
