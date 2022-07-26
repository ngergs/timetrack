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
 - begin: Starts a new time tracking session.
 - begin: Stops a previously started session.
 - status <YYYY-MM-dd,optional>: Prints the time tracking session status, the total working hour balance as well as the hours worked for the given date. If no date is entered for the optional input the current day is evaluated.

Further options:
```
Timetrack is a simple task manager that tracks the time spent working on a single project/job per day.
It automatically tracks the accumulated overtime. The daily working minutes are by default set to 480minutes=8hours.
To adjust this just define an alias, e.g. timetrack="timetrack -w 420" for 7hours per day.

Usage:
  timetrack [command]

Available Commands:
  begin       Starts a new timetrack sessions
  completion  Generate the autocompletion script for the specified shell
  end         Stops the current timetrack session
  help        Help about any command
  status      Prints the status of today's timesheet or from the specified date
  version     Prints the program version

Flags:
      --debug                     Set the log level to debug
  -h, --help                      help for timetrack
      --jsonlog                   Activates zerolog plain json output for the logs
  -d, --timetrack-folder string   Sets the folder where the timetrack time sheets are stored (default "~/.timetrack/")
  -w, --working-hours int         Daily working hours used to compute the time balance (default 800)

Use "timetrack [command] --help" for more information about a command.
```
