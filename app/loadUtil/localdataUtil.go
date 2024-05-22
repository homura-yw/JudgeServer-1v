package loadutil

import (
	"fmt"
	"io"
	"os"
)

type LocalDataUtil struct {
	path string
}

func (ld *LocalDataUtil) LoadToFile(src, dst string) error {
	bts := []byte(src)
	if bts[0] != '/' {
		src = "/" + src
	}
	src = ld.path + src

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func LocalDataUtilFactory(path string) (*LocalDataUtil, error) {
	return &LocalDataUtil{path: path}, nil
}
