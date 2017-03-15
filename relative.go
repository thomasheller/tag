package main

import (
	"log"
	"path/filepath"
	"strings"

	s "github.com/thomasheller/sortedset"
)

type RelativeTags struct {
	base Tags
	root string
	wd   string
}

func NewRelativeTags(base Tags, root string, wd string) *RelativeTags {
	return &RelativeTags{base: base, root: root, wd: wd}
}

func (t RelativeTags) Add(tag string, files ...string) {
	relFiles := t.rel(files)
	t.base.Add(tag, relFiles...)
}

func (t RelativeTags) Del(tag string, files ...string) {
	relFiles := t.rel(files)
	t.base.Del(tag, relFiles...)
}

func (t RelativeTags) Find(tag string) map[string]s.SortedSet {
	result := make(map[string]s.SortedSet)

	for file, tags := range t.base.Find(tag) {

		abs := filepath.Join(t.root, file)

		// filter files not in the current directory:

		if strings.HasPrefix(abs, t.wd) {
			// make paths relative:

			rel, err := filepath.Rel(t.wd, abs)
			if err != nil {
				log.Fatalf("Error getting relative path: %v", err)
			}

			result[rel] = tags
		}
	}

	return result
}

func (t RelativeTags) Untagged() []string {
	result := []string{}
	for _, file := range t.base.Untagged() {

		abs := filepath.Join(t.root, file)

		// filter files not in the current directory:

		if strings.HasPrefix(abs, t.wd) {
			// make paths relative:

			rel, err := filepath.Rel(t.wd, abs)
			if err != nil {
				log.Fatalf("Error getting relative path: %v", err)
			}

			result = append(result, rel)
		}
	}

	return result
}

func (t RelativeTags) List(prefix string) []string {
	rel, err := filepath.Rel(t.root, t.wd)
	if err != nil {
		log.Fatalf("Error getting relative path: %v", err)
	}

	return t.base.List(rel)
}

func (t RelativeTags) Dump() []string {
	if t.root == t.wd {
		return t.base.Dump()
	}

	result := []string{}

	rel, err := filepath.Rel(t.root, t.wd)
	if err != nil {
		log.Fatalf("Error getting relative path: %v", err)
	}

	for _, line := range t.base.Dump() {
		if strings.HasPrefix(line, rel) {
			result = append(result, line)
		}
	}

	return result
}

func (t RelativeTags) rel(files []string) []string {
	// input: foo.txt (which is /a/b/foo.txt)
	// output: b/foo.txt
	result := []string{}
	for _, file := range files {
		abs := filepath.Join(t.wd, file)

		rel, err := filepath.Rel(t.root, abs)
		if err != nil {
			log.Fatalf("Error getting relative path: %v", err)
		}

		result = append(result, rel)
	}
	return result
}
