package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type fileSystem interface {
	WriteOpen(filename string) error
	Fprintln(line string)
	FlushClose() error

	ReadOpen(filename string) error
	Scan() bool
	Text() string
	Err() error
	Close() error

	FileExists(filename string) bool

	Getwd() (string, error)
}

// osFileSystem is a simple wrapper around the file system, so we can
// mock it out when testing.
type osFileSystem struct {
	file *os.File
	w    *bufio.Writer
	sc   *bufio.Scanner
}

func (fs *osFileSystem) ReadOpen(filename string) error {
	if fs.file != nil {
		panic("Can't open another file at the same time!")
	}

	var err error
	fs.file, err = os.Open(filename)

	if err != nil {
		return err
	}

	fs.sc = bufio.NewScanner(fs.file)

	return nil
}

func (fs *osFileSystem) Scan() bool {
	if fs.sc == nil {
		panic("Scanner not available")
	}

	return fs.sc.Scan()
}

func (fs *osFileSystem) Text() string {
	if fs.sc == nil {
		panic("Scanner not available")
	}

	return fs.sc.Text()
}

func (fs *osFileSystem) Err() error {
	if fs.sc == nil {
		panic("Scanner not available")
	}

	return fs.sc.Err()
}

func (fs *osFileSystem) Close() error {
	if fs.file == nil {
		panic("Can't close yet, no open file!")
	}

	defer func() {
		fs.file = nil
		fs.sc = nil
	}()

	return fs.file.Close()
}

func (fs *osFileSystem) WriteOpen(filename string) error {
	if fs.file != nil {
		panic("Can't open another file at the same time!")
	}

	var err error
	fs.file, err = os.Create(filename)

	if err != nil {
		return err
	}

	fs.w = bufio.NewWriter(fs.file)

	return nil
}

func (fs *osFileSystem) Fprintln(line string) {
	if fs.file == nil {
		panic("Can't write line before opening a file!")
	}

	fmt.Fprintln(fs.w, line)
}

func (fs *osFileSystem) FlushClose() error {
	if fs.file == nil {
		panic("Can't flush or close yet, no open file!")
	}

	defer func() {
		fs.file = nil
		fs.w = nil
	}()

	err := fs.w.Flush()
	if err != nil {
		return err
	}

	fs.file.Close()

	return nil
}

func (fs *osFileSystem) FileExists(filename string) bool {
	file, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Fatalf("Error checking if file %s exists: %v", filename, err)
	}
	if file.IsDir() {
		return false
	}
	return true
}

func (fs *osFileSystem) Getwd() (string, error) {
	return os.Getwd()
}
