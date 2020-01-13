package system

import (
	"io/ioutil"
	"os"
)

func TemporaryFile(directory string, outputPattern string) (*os.File, error) {
	file, err := ioutil.TempFile(directory, outputPattern)
	if err != nil {
		return nil, err
	}

	return file, nil
}
