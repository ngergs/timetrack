package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestReadWrite(t *testing.T) {
	// create temp folder if non existing
	tmpDir := "test/tmp"
	err := os.MkdirAll(tmpDir, 0755)
	assert.Nil(t, err)

	filename := path.Join(tmpDir, "readwrite.dat")
	dummyData := "testmessage"
	err = write(filename, &dummyData)
	assert.Nil(t, err)
	defer os.Remove(path.Join(resolvedFolder, filename))
	readData, err := read[string](filename)
	assert.Nil(t, err)
	assert.Equal(t, dummyData, *readData)
}
