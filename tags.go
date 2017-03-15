package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	s "github.com/thomasheller/sortedset"
)

type Tags interface {
	Add(tag string, files ...string)
	Del(tag string, files ...string)
	Find(tag string) map[string]s.SortedSet
	Untagged() []string
	List(prefix string) []string
	Dump() []string
}

type BaseTags struct {
	db   Db
	w    Walker
	fs   fileSystem
	root string
}

// NewBaseTags returns a new Tags object. Assumes that Db is in loaded
// state.
func NewBaseTags(db Db, w Walker, fs fileSystem, root string) Tags {
	return BaseTags{db: db, w: w, fs: fs, root: root}
}

func (t BaseTags) Add(tag string, files ...string) {
	for _, file := range files {
		if !t.fs.FileExists(filepath.Join(t.root, file)) {
			fmt.Printf("skipping non-existent file: %s\n", file)
			continue
		}

		t.db.add(file, tag)
	}

	t.db.save()
}

func (t BaseTags) Del(tag string, files ...string) {
	for _, file := range files {
		if !t.fs.FileExists(filepath.Join(t.root, file)) {
			fmt.Printf("skipping non-existent file: %s\n", file)
			continue
		}

		t.db.remove(file, tag)
	}

	t.db.save()
}

func (t BaseTags) Find(tag string) map[string]s.SortedSet {
	matches := make(map[string]s.SortedSet)

	for filename, tags := range t.db.list() {
		if tags.Contains(tag) {
			// match := fmt.Sprintf("%s:%s", filename, tags.String())
			// matches = append(matches, match)

			matches[filename] = tags
		}
	}

	// sort.Strings(matches)

	return matches
}

func (t BaseTags) Untagged() []string {
	untagged := []string{}

WalkLoop:
	for _, filename := range t.w.Walk(t.root) {
		for fileWithTag := range t.db.list() {
			if filename == fileWithTag {
				continue WalkLoop
			}
		}
		untagged = append(untagged, filename)
	}

	sort.Strings(untagged)

	return untagged
}

func (t BaseTags) List(prefix string) []string {
	list := s.New([]string{})

	for filename, tags := range t.db.list() {
		if prefix == "." || strings.HasPrefix(filename, prefix) {
			for _, tag := range tags.Slice() {
				list.Add(tag)
			}
		}
	}

	result := []string{}

	for _, tag := range list.Slice() {
		result = append(result, tag)
	}

	return result
}

func (t BaseTags) Dump() []string {
	return t.db.dump()
}
