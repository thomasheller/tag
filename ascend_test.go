package main

import (
	"testing"
)

func TestAscend(t *testing.T) {
	fs := newMemoryFileSystem("/home/foo/bar")
	fs.Touch("/home/foo/my.txt")
	a := newAscender(fs)
	file, path, err := a.ascend("my.txt")
	if err != nil {
		t.Fatal(err)
	}
	if file != "/home/foo/my.txt" {
		t.Fatalf("my.txt at unexpected location %s", file)
	}
	if path != "/home/foo" {
		t.Fatalf("my.txt in unexpected directory %s", path)
	}
}
