package main

import (
	"fmt"
	"github.com/ngergs/timetrack/v2/internal/constants"
	"github.com/ngergs/timetrack/v2/internal/io"
	"github.com/ngergs/timetrack/v2/sheet"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
	"time"
)

const timeTrackFolderFlag = "timetrack-folder"
const workingMinutesFlag = "working-minutes"
const debugFlag = "debug"
const jsonlogFlag = "jsonlog"

// set via flags in init(), constant afterwards
var timeTrackFolder string
var workingMinutes int

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "timetrack",
	Short: "Tracks the time spent working during the day",
	Long: `Timetrack is a simple task manager that tracks the time spent working on a single project/job per day.
It automatically tracks the accumulated overtime. The daily working minutes are by default set to 480minutes=8hours.
To adjust this just define an alias, e.g. timetrack="timetrack -w 420" for 7hours per day.`,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().StringP(timeTrackFolderFlag, "d", "~/.timetrack/", "Sets the folder where the timetrack time sheets are stored")
	rootCmd.PersistentFlags().IntP(workingMinutesFlag, "w", 480, "Daily working time in minutes, used to compute the time balance")
	rootCmd.PersistentFlags().Bool(debugFlag, false, "Set the log level to debug")
	rootCmd.PersistentFlags().Bool(jsonlogFlag, false, "Activates zerolog plain json output for the logs")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		err := initLogging(cmd)
		if err != nil {
			return err
		}
		return initFlagsAndFolder(cmd)
	}
}

func initLogging(cmd *cobra.Command) error {
	jsonlogActive, err := cmd.Flags().GetBool(jsonlogFlag)
	if err != nil {
		return fmt.Errorf("could not parse jsonlog logging setting: %w", err)
	}
	if !jsonlogActive {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	debugActive, err := cmd.Flags().GetBool(debugFlag)
	if err != nil {
		return fmt.Errorf("could not parse debug setting: %w", err)
	}
	if debugActive {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	return nil
}

func initFlagsAndFolder(cmd *cobra.Command) (err error) {
	timeTrackFolder, err = cmd.Flags().GetString(timeTrackFolderFlag)
	if err != nil {
		return fmt.Errorf("could not determine timetrack sheet folder: %w", err)
	}
	workingMinutes, err = cmd.Flags().GetInt(workingMinutesFlag)
	if err != nil {
		return fmt.Errorf("could not determine daily working hours: %w", err)
	}
	if strings.Contains(timeTrackFolder, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to obtain user home dir: %w", err)
		}
		timeTrackFolder = strings.ReplaceAll(timeTrackFolder, "~", home)
	}
	err = os.MkdirAll(timeTrackFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create timekeep timestamps folder: %w", err)
	}
	return nil
}

func loadSheetAndDo(date *time.Time, writeChanges bool, action func(*sheet.Timesheet) error) error {
	file, timesheet, err := getCurrentTimesheet(date, writeChanges)
	if err != nil {
		return err
	}
	defer io.Close(file)

	err = action(timesheet)
	if err != nil {
		return err
	}

	if writeChanges {
		err = io.Write(file, timesheet)
		if err != nil {
			return fmt.Errorf("failed to save modified timesheet: %w", err)
		}
	}
	return nil
}

func getCurrentTimesheet(date *time.Time, writeChanges bool) (*os.File, *sheet.Timesheet, error) {
	file, err := sheet.GetLastSavedFile(timeTrackFolder, date, writeChanges)
	if err != nil {
		return nil, nil, err
	}
	var timesheet *sheet.Timesheet
	if file != nil {
		timesheet, err = sheet.GetLastTimesheet(file)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load timesheet: %w", err)
		}
	} else {
		// no timesheet present, create a new one and prepare file on disk for save
		log.Debug().Msg("No timesheet found, opening a new one")
		timesheet = sheet.New()
		if writeChanges {
			file, err = io.OpenFile(path.Join(timeTrackFolder, sheet.CurrentDayString), writeChanges)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to open new timesheet: %w", err)
			}
		}
	}

	// timesheet is from last day, rotate and create a new one
	if len(timesheet.Slices) > 0 && timesheet.Slices[0].Start.Format(constants.ReferenceFormat) != sheet.CurrentDayString {
		oldTimesheet := timesheet
		timesheet, err = sheet.PrepareNewFromOldSheet(timesheet, workingMinutes)
		// finish oldTimesheet timesheet file and open a new one
		if writeChanges {
			newFile, err := io.OpenFile(path.Join(timeTrackFolder, sheet.CurrentDayString), true)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to open new timesheet for today: %w", err)
			}
			// only write changes to the previous file when the new one is successfully opened
			err = io.Write(file, oldTimesheet)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to write closing statement to oldTimesheet timesheet: %w", err)
			}
			io.Close(file)
			file = newFile
		}
	}
	return file, timesheet, nil
}
