package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

var filePattern, _ = regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2}\\.json")

func getLastSavedTimesheet(dir string) (sheet *timesheet, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// entries are sorted by filename, so start from the end
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if filePattern.MatchString(entry.Name()) &&
			entry.Type().IsRegular() {
			var sheet *timesheet // goland limitation, temporary
			sheet, err := read[timesheet](path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			err = sheet.Validate()
			if err != nil {
				return nil, err
			}
			return sheet, nil
		}
	}
	return nil, nil
}

func read[T any](path string) (*T, error) {
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

func write[T any](path string, content *T) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0755)
}
