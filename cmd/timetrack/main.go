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

const timeTrackFolder = "timetrack-folder"
const workingMinutes = "working-minutes"
const debug = "debug"
const jsonlog = "jsonlog"

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
	rootCmd.PersistentFlags().StringP(timeTrackFolder, "d", "~/.timetrack/", "Sets the folder where the timetrack time sheets are stored")
	rootCmd.PersistentFlags().IntP(workingMinutes, "w", 480, "Daily working time in minutes, used to compute the time balance")
	rootCmd.PersistentFlags().Bool(debug, false, "Set the log level to debug")
	rootCmd.PersistentFlags().Bool(jsonlog, false, "Activates zerolog plain json output for the logs")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		jsonlogActive, err := cmd.Flags().GetBool(jsonlog)
		if err != nil {
			return fmt.Errorf("could not parse jsonlog logging setting: %w", err)
		}
		if !jsonlogActive {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		debugActive, err := cmd.Flags().GetBool(debug)
		if err != nil {
			return fmt.Errorf("could not parse debug setting: %w", err)
		}
		if debugActive {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		resolvedFolder, err := cmd.Flags().GetString(timeTrackFolder)
		if err != nil {
			return fmt.Errorf("could not parse timetrack-folder setting: %w", err)
		}
		if strings.Contains(resolvedFolder, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to obtain user home dir: %w", err)
			}
			resolvedFolder = strings.ReplaceAll(resolvedFolder, "~", home)
		}
		err = os.MkdirAll(resolvedFolder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create timekeep timestamps folder: %w", err)
		}
		err = rootCmd.Flags().Set(timeTrackFolder, resolvedFolder)
		if err != nil {
			return fmt.Errorf("failed to set resolved timetrack-folder flag: %w", err)
		}
		return nil
	}
}

func loadSheetAndDo(cmd *cobra.Command, date *time.Time, writeChanges bool, action func(*sheet.Timesheet) error) error {
	folder, err := cmd.Flags().GetString(timeTrackFolder)
	if err != nil {
		return fmt.Errorf("could not determine timetrack sheet folder: %w", err)
	}
	workingMinutes, err := cmd.Flags().GetInt(workingMinutes)
	if err != nil {
		return fmt.Errorf("could not determine daily working hours: %w", err)
	}

	file, err := sheet.GetLastSavedFile(folder, date, writeChanges)
	if err != nil {
		return err
	}
	if file != nil {
		defer io.Close(file)
	}

	current, err := sheet.GetLastTimesheet(file)
	if err != nil {
		return fmt.Errorf("failed to load timesheet: %w", err)
	}

	// prepare a save file when the current file lives purely in memory and writing has been requested
	if writeChanges && file == nil {
		file, err = io.OpenFile(path.Join(folder, sheet.CurrentDayString), writeChanges)
		if err != nil {
			return fmt.Errorf("failed to open new timesheet: %w", err)
		}
		defer io.Close(file)
	}

	if len(current.Slices) > 0 && current.Slices[0].Start.Format(constants.ReferenceFormat) != sheet.CurrentDayString {
		old := current
		current, err = sheet.PrepareNewFromOldSheet(current, workingMinutes)
		// finish old timesheet file and open a new one
		if writeChanges {
			newFile, err := io.OpenFile(path.Join(folder, sheet.CurrentDayString), true)
			if err != nil {
				return fmt.Errorf("failed to open new timesheet for today: %w", err)
			}
			defer io.Close(newFile)
			// only write changes to the previous file when the new one is successfully opened
			err = io.Write(file, old)
			if err != nil {
				return fmt.Errorf("failed to write closing statement to old timesheet: %w", err)
			}
			file = newFile
		}
	}

	err = action(current)
	if err != nil {
		return err
	}

	if writeChanges {
		err = io.Write(file, current)
		if err != nil {
			return fmt.Errorf("failed to save modified timesheet: %w", err)
		}
	}
	return nil
}
