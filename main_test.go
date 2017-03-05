package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// TODO: test in subdirectories

func setupTags() (Tags, *memoryFileSystem) {
	fs := newMemoryFileSystem("/home/foo")

	w := newMemoryFileSystemWalker(*fs)

	db := newDb("tags.dat", fs)
	db.load()

	tags := NewBaseTags(*db, w, fs, "/home/foo")

	return tags, fs
}

func ExampleAdd() {
	tags, fs := setupTags()
	fs.Touch("a.txt")
	runCommand(tags, "add", "foo", "a.txt")
	runCommand(tags, "list")
	// Output:
	// foo
}

func ExampleAddNonExistent() {
	tags, fs := setupTags()
	fs.Touch("a.txt")
	runCommand(tags, "add", "foo", "b.txt")
	runCommand(tags, "list")
	// Output:
	// skipping non-existent file: b.txt
}

func ExampleDel() {
	tags, fs := setupTags()
	fs.Touch("a.txt")
	runCommand(tags, "add", "foo", "a.txt")
	runCommand(tags, "del", "foo", "a.txt")
	runCommand(tags, "list")
	// Output:
}

func ExampleList() {
	tags, fs := setupTags()
	fs.Touch("a.txt")
	fs.Touch("b.txt")
	runCommand(tags, "add", "foo", "a.txt")
	runCommand(tags, "add", "bar", "b.txt")
	runCommand(tags, "list")
	// Output:
	// bar
	// foo
}

func ExampleFind() {
	tags, fs := setupTags()
	fs.Touch("a.txt")
	fs.Touch("b.txt")
	fs.Touch("c.txt")
	runCommand(tags, "add", "foo", "b.txt", "c.txt", "a.txt")
	runCommand(tags, "add", "bar", "c.txt")
	runCommand(tags, "find", "foo")
	// Output:
	// a.txt:foo
	// b.txt:foo
	// c.txt:bar,foo
}

func ExampleUntagged() {
	tags, fs := setupTags()
	fs.Touch("a.txt")
	fs.Touch("b.txt")
	fs.Touch("c.txt")
	runCommand(tags, "add", "foo", "a.txt")
	runCommand(tags, "untagged")
	// Output:
	// b.txt
	// c.txt
	// tags.dat
}

func ExampleDump() {
	tags, fs := setupTags()
	fs.Touch("a.txt")
	fs.Touch("b.txt")
	runCommand(tags, "add", "bar", "a.txt", "b.txt")
	runCommand(tags, "add", "foo", "a.txt")
	runCommand(tags, "dump")
	// Output:
	// a.txt:bar,foo
	// b.txt:bar
}

// memoryFileSystemWalker walks over a simulated file system.
type memoryFileSystemWalker struct {
	fs memoryFileSystem
}

func newMemoryFileSystemWalker(mfs memoryFileSystem) Walker {
	return memoryFileSystemWalker{fs: mfs}
}

// Walk returns all files from the memoryFileSystem. The root
// parameter is ignored.
func (w memoryFileSystemWalker) Walk(root string) []string {
	files := make([]string, 0)
	for _, filename := range w.fs.Files() {
		files = append(files, filename)
	}
	return files
}

// memoryFileSystem provides a simulated file system for testing,
// without actually writing files to disk.
type memoryFileSystem struct {
	cwd         string
	files       map[string]string
	currentFile string
	sc          *bufio.Scanner
}

func newMemoryFileSystem(cwd string) *memoryFileSystem {
	return &memoryFileSystem{cwd: cwd, files: make(map[string]string)}
}

func (fs *memoryFileSystem) Files() []string {
	files := make([]string, 0)
	for filename := range fs.files {
		relative := filename[len(fs.cwd)+1:]
		files = append(files, relative)
	}
	return files
}

func (fs *memoryFileSystem) WriteOpen(filename string) error {
	if fs.currentFile != "" {
		panic("Can't open file, another file is open!")
	}
	fs.currentFile = fs.abs(filename)
	fs.files[fs.currentFile] = ""
	return nil
}

func (fs *memoryFileSystem) Fprintln(line string) {
	fs.files[fs.currentFile] = fs.files[fs.currentFile] + line + "\n"
}

func (fs *memoryFileSystem) FlushClose() error {
	fs.currentFile = ""
	return nil
}

func (fs *memoryFileSystem) ReadOpen(filename string) error {
	if fs.currentFile != "" {
		panic("Can't open file, another file is open!")
	}
	fs.currentFile = fs.abs(filename)
	fs.sc = bufio.NewScanner(strings.NewReader(fs.files[fs.currentFile]))
	return nil
}

func (fs *memoryFileSystem) Scan() bool {
	return fs.sc.Scan()
}

func (fs *memoryFileSystem) Text() string {
	return fs.sc.Text()
}

func (fs *memoryFileSystem) Err() error {
	return fs.sc.Err()
}

func (fs *memoryFileSystem) Close() error {
	fs.currentFile = ""
	return nil
}

func (fs *memoryFileSystem) FileExists(filename string) bool {
	_, ok := fs.files[fs.abs(filename)]
	return ok
}

func (fs *memoryFileSystem) Getwd() (string, error) {
	return fs.cwd, nil
}

func (fs *memoryFileSystem) Touch(filename string) {
	fs.files[fs.abs(filename)] = ""
}

func (fs *memoryFileSystem) abs(path string) string {
	if fs.isAbs(path) {
		return path
	} else {
		return filepath.Join(fs.cwd, path)
	}
}

func (fs *memoryFileSystem) isAbs(path string) bool {
	return len(path) > 0 && path[0] == os.PathSeparator
}
