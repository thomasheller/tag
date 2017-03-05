package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type ascender struct {
	fs fileSystem
}

func newAscender(fs fileSystem) *ascender {
	return &ascender{fs: fs}
}

func (a *ascender) ascend(filename string) (string, string, error) {
	cwd, err := a.fs.Getwd()
	if err != nil {
		return "", "", err
	}

	for dir := cwd; dir != string(os.PathSeparator); dir = path.Dir(dir) {
		path := path.Join(dir, filename)

		if a.fs.FileExists(path) {
			return path, filepath.Dir(path), nil
		}
	}

	return "", "", fmt.Errorf("%s not found in %s (or any parent directory)", filename, cwd)
}
