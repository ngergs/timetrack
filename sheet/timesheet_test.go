package sheet

import (
	"math"
	"testing"
	"time"

	"github.com/ngergs/timetrack/v2/sheet/states"
	"github.com/stretchr/testify/assert"
)

func TestStartStopSession(t *testing.T) {
	sheet := New()
	// empty sheets are ok for new initialization
	validateExpectOk(t, sheet)
	assert.Equal(t, 0, len(sheet.Slices))
	assert.Equal(t, states.Closed, sheet.GetState())
	startSessionExpectOk(t, sheet, 0)

	// starting a second session is not possible while on is running
	err := sheet.BeginSession()
	validateExpectOk(t, sheet)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(sheet.Slices))
	assert.Equal(t, states.Open, sheet.GetState())

	// stopping the session should work
	err = sheet.EndSession()
	validateExpectOk(t, sheet)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(sheet.Slices))
	assert.Equal(t, states.Closed, sheet.GetState())
	assert.NotNil(t, sheet.Slices[0].End)
	assert.Less(t, time.Since(*sheet.Slices[0].End).Milliseconds(), int64(100))

	// double stop should fail
	err = sheet.EndSession()
	validateExpectOk(t, sheet)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(sheet.Slices))
	assert.Equal(t, states.Closed, sheet.GetState())

	// starting a second session should now work
	startSessionExpectOk(t, sheet, 1)
}

func TestBalance(t *testing.T) {
	timeDifference := time.Duration(42)*time.Second + time.Duration(12)*time.Minute + time.Duration(7)*time.Hour
	startTime := time.Now().Add(-timeDifference)
	endTime := time.Now()
	// update timeDifference to take nanoSeconds drift into account
	timeDifference = endTime.Sub(startTime)

	sheet := New()
	sheet.Slices = append(sheet.Slices, Timeslice{
		Start: &startTime,
		End:   &endTime,
	})
	assert.Equal(t, int(math.Floor(timeDifference.Minutes())), sheet.GetTodayBalance())
}

func startSessionExpectOk(t *testing.T, sheet *Timesheet, expectedSessionCount int) {
	err := sheet.BeginSession()
	validateExpectOk(t, sheet)
	assert.Nil(t, err)
	assert.Equal(t, expectedSessionCount+1, len(sheet.Slices))
	assert.Equal(t, states.Open, sheet.GetState())
	assert.NotNil(t, sheet.Slices[expectedSessionCount].Start)
	assert.Nil(t, sheet.Slices[expectedSessionCount].End)
	assert.Less(t, time.Since(*sheet.Slices[expectedSessionCount].Start).Milliseconds(), int64(100))
}

func validateExpectOk(t *testing.T, sheet *Timesheet) {
	err := sheet.Validate()
	assert.Nil(t, err)
}
