package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	s "github.com/thomasheller/sortedset"
)

func main() {
	if len(os.Args) == 1 {
		usage()
	}

	op := os.Args[1]
	params := os.Args[2:]

	fs := &osFileSystem{}

	a := newAscender(fs)
	filename, root, err := a.ascend("tags.dat")

	if err != nil {
		log.Fatalf("Unable to find tags database: %v", err)
	}

	w := &FilePathWalker{}
	db := newDb(filename, fs)
	db.load()
	base := NewBaseTags(*db, w, fs, root)
	wd, err := fs.Getwd()
	if err != nil {
		log.Fatalf("Unable to find current working directory: %v", err)
	}
	// TODO: relative tags is require only if root != wd
	tags := NewRelativeTags(base, root, wd)

	if !runCommand(tags, op, params...) {
		usage()
	}
}

// runCommand runs the specified operation with corresponding
// parameters, returns false if op/parameters couldn't be parsed and a
// usage mesage should be displayed instead.
func runCommand(tags Tags, op string, params ...string) bool {
	switch op {
	case "add":
		if len(params) < 2 {
			return false
		}
		tags.Add(params[0], params[1:]...)
	case "del":
		if len(params) < 2 {
			return false
		}
		tags.Del(params[0], params[1:]...)
	case "find":
		if len(params) != 1 {
			return false
		}
		printSortedMap(tags.Find(params[0]))
	case "untagged":
		if len(params) != 0 {
			return false
		}
		printLines(tags.Untagged())
	case "list":
		if len(params) != 0 {
			return false
		}
		printLines(tags.List("."))
	case "dump":
		if len(params) != 0 {
			return false
		}
		printLines(tags.Dump())
	default:
		return false
	}

	return true
}

func printSortedMap(m map[string]s.SortedSet) {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		tags := m[key]
		line := fmt.Sprintf("%s:%s", key, tags.String())
		fmt.Println(line)
	}
}

func printLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

// usage prints the usage message.
func usage() {
	fmt.Println("Usage:")
	fmt.Println("  tag [command] [parameter...]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println()
	fmt.Println("  tag add   [tag] [file...]   add tag to file(s)")
	fmt.Println("  tag del   [tag] [file...]   remove tag from file(s)")
	fmt.Println("  tag find  [tag]             find files with tag in current directory (recursive)")
	fmt.Println("  tag list                    list all tags used in current directory (recursive)")
	fmt.Println("  tag untagged                list untagged files in current directory (recursive)")
	fmt.Println("  tag dump                    dump entries from tags.dat file (debugging)")
	os.Exit(0)
}
