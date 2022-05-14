# Timetrack
A time-tracking utility following the KISS principle. This utility is meant to track working hours.
The config options allow to set the daily working hours agreed upon as well as an alternative storage place
for the tracked time sheets (default is ~/.timetrack).

## Installation
### Binary
Each release has pre-compiled binaries for linux/osx/windows attached. You can just download and use them :)

### Building from source
Clone this repository and run (with go v1.18+ installed)
```bash
go build
```

## Usage
Three operation modes are supported:
 - start: Starts a new time tracking session.
 - stop: Stops a previously started session.
 - status: Prints the time tracking session status, the total working hour balance as well as the hours worked today.
 - sheet YYYY-MM-dd: Prints the status command from the retrospective perspective of the specified day (a time sheet has to exist for that day).

Further options:
```
Usage: ./timetrack {options} (start|stop|status|sheet {date YYYY-MM-DD})
Options:
  -compact
    	whether the status should be printed in a compact format omitting timesheet details
  -debug
    	Log debug level
  -folder string
    	folder in which the timeekeep time slice are saved (default "~/.timetrack/")
  -help
    	Prints the help.
  -pretty
    	Activates zerolog pretty logging (default true)
  -version
        Prints the version of this program
  -working-minutes int
    	daily working minutes to update the time balance (default 480)
```
