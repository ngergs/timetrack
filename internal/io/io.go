package io

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

// OpenFile opens a file for reading and if the write argument is set also for writing
func OpenFile(path string, write bool) (*os.File, error) {
	var mode int
	if write {
		mode |= os.O_CREATE | os.O_RDWR
	} else {
		mode |= os.O_RDONLY
	}
	return os.OpenFile(path, mode, 0755)
}

func Read[T any](r io.Reader) (*T, error) {
	data, err := io.ReadAll(r)
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

func Write[T any](f *os.File, content *T) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	var seeker int64
	for len(data) != 0 {
		n, err := f.WriteAt(data, seeker)
		if err != nil {
			return err
		}
		seeker += int64(n)
		data = data[n:]
	}
	return nil
}

// Close closes the io.Closer and logs error as warnings.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Warn().Err(err).Msg("Could not close file")
	}
}
