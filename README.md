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

Further options:
```
Usage: ./timetrack {options} (start|stop|status)
Options:
  -debug
        Log debug level
  -folder string
        folder in which the timeekeep time slice are saved (default "~/.timetrack/")
  -help
        Prints the help.
  -pretty
        Activates zerolog pretty logging (default true)
  -working-minutes int
        daily working minutes to update the time balance (default 480)
```
