package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

// exists returns whether the given file or directory exists
func isFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// travers up directories until we find a main.go file
func findRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	for i := 0; i < 3; i++ {
		if ok, _ := isFileExists(path.Join(dir, "main.go")); ok {
			return dir
		}
		dir = path.Dir(dir)
	}

	panic(fmt.Errorf("fatal error config file: %w", err))
}
