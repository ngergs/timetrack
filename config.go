package main

import (
	"flag"
	"fmt"
	"github.com/ngergs/timetrack/modes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

var version = "snapshot"
var compactPrint = flag.Bool("compact", false, "whether the status should be printed in a compact format omitting timesheet details")
var folder = flag.String("folder", "~/.timetrack/", "folder in which the timeekeep time slice are saved")
var dailyWorkingMinutes = flag.Int("working-minutes", 480, "daily working minutes to update the time balance")
var prettyLogging = flag.Bool("pretty", true, "Activates zerolog pretty logging")
var debugLogging = flag.Bool("debug", false, "Log debug level")
var help = flag.Bool("help", false, "Prints the help.")
var printVersion = flag.Bool("version", false, "Prints the version of this program")

var mode modes.Mode
var resolvedFolder string

func getArgs(expectedLength int) []string {
	args := flag.Args()
	if len(args) != expectedLength {
		log.Error().Msgf("Unexpected number of arguments: %d, expected %d", len(args), expectedLength)
		flag.Usage()
		os.Exit(1)
	}
	return args
}
func getArgsMinLength(minLength int) []string {
	args := flag.Args()
	if len(args) < minLength {
		log.Error().Msgf("Unexpected number of arguments: %d, expected more or equal than %d", len(args), minLength)
		flag.Usage()
		os.Exit(1)
	}
	return args
}

func readConfig() {
	if *prettyLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s {options} (start|stop|status|sheet {date YYYY-MM-DD})\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *printVersion {
		fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", version)
		os.Exit(0)
	}
	if *debugLogging {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	args := getArgsMinLength(1)
	var ok bool
	mode, ok = modes.Parse(args[0])
	if !ok {
		log.Error().Msgf("Unsupported operating modes: %s", args[0])
		flag.Usage()
		os.Exit(1)
	}

	resolvedFolder = *folder
	if strings.Contains(*folder, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to obtain user home dir")
		}
		resolvedFolder = strings.ReplaceAll(resolvedFolder, "~", home)
	}
	err := os.MkdirAll(resolvedFolder, 0755)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create timekeep timestamps folder")
	}
}
