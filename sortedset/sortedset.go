package sortedset

import (
	"sort"
	"strings"
)

// SortedSet represents an alphabetically sorted set of strings.
type SortedSet struct {
	data []string
}

// New builds a new set from the given elements.
func New(elements []string) SortedSet {
	l := SortedSet{data: make([]string, 0)}
	for _, e := range elements {
		if !l.Contains(e) {
			l.data = append(l.data, e)
		}
	}
	sort.Strings(l.data)
	return l
}

// Add adds an element to the set. If it already exists, nothing
// happens, as set elements are unique.
func (l *SortedSet) Add(element string) {
	if len(l.data) == 0 {
		l.data = []string{element}
		return
	}

	if l.Contains(element) {
		return
	}

	l.data = append(l.data, element)
	sort.Strings(l.data)
}

// Delete removes an element from the set. If it doesn't exist,
// nothing happens. Returns true if the set is now empty.
func (l *SortedSet) Delete(element string) (empty bool) {
	if len(l.data) == 0 {
		return true
	}

	if len(l.data) == 1 && l.data[0] == element {
		l.data = make([]string, 0)
		return true
	}

	for i, e := range l.data {
		if e == element {
			l.data = append(l.data[:i], l.data[i+1:]...)
			return false
		}
	}

	return false
}

// String returns all elements separated by commas.
func (l *SortedSet) String() string {
	return strings.Join(l.data, ",")
}

// Slice returns all elements as a slice.
func (l *SortedSet) Slice() []string {
	result := make([]string, len(l.data))
	copy(result, l.data)
	return result
}

// Contains reports if the set contains the given element.
func (l *SortedSet) Contains(search string) bool {
	for _, e := range l.data {
		if search == e {
			return true
		}
	}
	return false
}

// ContainsAll reports if the set contains all elements given in
// search.
func (l *SortedSet) ContainsAll(search []string) bool {
SearchLoop:
	for _, s := range search {
		for _, e := range l.data {
			if s == e {
				continue SearchLoop
			}
		}
		return false
	}
	return true
}

// ContainsAll reports if the set contains any of the elements given
// in search.
func (l *SortedSet) ContainsAny(search []string) bool {
	for _, s := range search {
		for _, e := range l.data {
			if s == e {
				return true
			}
		}
	}
	return false
}
