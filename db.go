package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	s "github.com/thomasheller/sortedset"
)

// Db implements a primitive key-value store.
type Db struct {
	filename string
	fs       fileSystem
	data     map[string]s.SortedSet
}

func newDb(filename string, fs fileSystem) *Db {
	return &Db{filename: filename, fs: fs}
}

func (db *Db) load() {
	db.data = make(map[string]s.SortedSet)

	if !db.fs.FileExists(db.filename) {
		return
	}

	err := db.fs.ReadOpen(db.filename)
	if err != nil {
		log.Fatal(err)
	}

	defer db.fs.Close()

	for db.fs.Scan() {
		line := db.fs.Text()
		parts := strings.Split(line, ":")
		k := parts[0]
		v := parts[1]
		values := strings.Split(v, ",")
		db.data[k] = s.New(values)
	}

	if err := db.fs.Err(); err != nil {
		log.Fatal(err)
	}
}

func (db *Db) save() {
	lines := make([]string, 0)

	for k, v := range db.data {
		line := fmt.Sprintf("%s:%s", k, v.String())
		lines = append(lines, line)
	}

	sort.Strings(lines)

	err := db.fs.WriteOpen(db.filename)
	if err != nil {
		panic(err)
	}

	for _, line := range lines {
		db.fs.Fprintln(line)
	}

	db.fs.FlushClose()
}

func (db *Db) add(key string, value string) {
	if strings.ContainsRune(key, '\n') {
		log.Fatalf("key \"%s\" must not contain newline", key)
	}
	if strings.ContainsRune(value, '\n') {
		log.Fatalf("value \"%s\" must not contain newline", value)
	}
	if strings.ContainsRune(key, ':') {
		log.Fatalf("key \"%s\" must not contain colon", key)
	}
	if strings.ContainsRune(value, ':') {
		log.Fatalf("value \"%s\" must not contain colon", value)
	}
	if strings.ContainsRune(value, ',') {
		log.Fatalf("value \"%s\" must not contain comma", value)
	}

	set, ok := db.data[key]

	if ok {
		set.Add(value)
		db.data[key] = set
	} else {
		db.data[key] = s.New([]string{value})
	}
}

func (db *Db) remove(key string, value string) {
	set, ok := db.data[key]
	if !ok {
		return
	}

	empty := set.Delete(value)

	if empty {
		delete(db.data, key)
	} else {
		db.data[key] = set
	}
}

func (db *Db) list() map[string]s.SortedSet {
	result := make(map[string]s.SortedSet)
	for k, v := range db.data {
		result[k] = v
	}
	return result
}

func (db *Db) dump() []string {
	err := db.fs.ReadOpen(db.filename)
	if err != nil {
		log.Fatal(err)
	}

	defer db.fs.Close()

	result := []string{}

	for db.fs.Scan() {
		line := db.fs.Text()
		result = append(result, line)
	}

	if err := db.fs.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}
