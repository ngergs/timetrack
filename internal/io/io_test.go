package io

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadWrite(t *testing.T) {
	// create temp folder if non existing
	tmpDir := "test/tmp"
	err := os.MkdirAll(tmpDir, 0755)
	assert.Nil(t, err)
	dummyData := "testmessage"
	file, err := ioutil.TempFile(tmpDir, "testFile")
	assert.Nil(t, err)
	defer os.Remove(file.Name())

	err = Write(file, &dummyData)
	assert.Nil(t, err)
	readData, err := Read[string](file)
	assert.Nil(t, err)
	assert.Equal(t, dummyData, *readData)
}
