package sortedset

import (
	"github.com/thomasheller/slicecmp"
	"testing"
)

func TestAdd(t *testing.T) {
	l := New([]string{"foo", "bar"})
	l.Add("ccc")

	if !slicecmp.Equal([]string{"bar", "ccc", "foo"}, l.Slice()) {
		t.Fatal("Unexpected")
	}
}

func TestAddTwice(t *testing.T) {
	l := New([]string{"foo", "bar"})
	l.Add("ccc")
	l.Add("ccc")

	if !slicecmp.Equal([]string{"bar", "ccc", "foo"}, l.Slice()) {
		t.Fatal("Unexpected")
	}
}

func TestDelete(t *testing.T) {
	l := New([]string{"foo", "bar"})
	l.Delete("foo")

	if !slicecmp.Equal([]string{"bar"}, l.Slice()) {
		t.Fatal("Unexpected")
	}
}

func TestDelete2(t *testing.T) {
	l := New([]string{"123", "34f0jf", "aaa"})
	l.Delete("123")

	if !slicecmp.Equal([]string{"34f0jf", "aaa"}, l.Slice()) {
		t.Fatal("Unexpected")
	}
}

func TestContainsAny(t *testing.T) {
	l := New([]string{"foo", "bar"})
	found := l.ContainsAny([]string{"foo", "zzz"})
	if !found {
		t.Fatal("Unexpected")
	}
}

func TestContainsAnyNotFound(t *testing.T) {
	l := New([]string{"foo", "bar"})
	found := l.ContainsAny([]string{"xxx", "zzz"})
	if found {
		t.Fatal("Unexpected")
	}
}

func TestContainsAll(t *testing.T) {
	l := New([]string{"foo", "bar"})
	found := l.ContainsAll([]string{"foo", "bar"})
	if !found {
		t.Fatal("Unexpected")
	}
}

func TestContainsAllNotFoundTooMany(t *testing.T) {
	l := New([]string{"foo", "bar"})
	found := l.ContainsAll([]string{"foo", "bar", "baz"})
	if found {
		t.Fatal("Unexpected")
	}
}

func TestContainsAllNotFoundTooFew(t *testing.T) {
	l := New([]string{"foo", "bar"})
	found := l.ContainsAll([]string{"foo"})
	if !found {
		t.Fatal("Unexpected")
	}
}

func TestContainsAllNotFoundPartial(t *testing.T) {
	l := New([]string{"foo", "bar"})
	found := l.ContainsAll([]string{"foo", "xxx"})
	if found {
		t.Fatal("Unexpected")
	}
}

func TestContainsAllNotFoundDifferent(t *testing.T) {
	l := New([]string{"foo", "bar"})
	found := l.ContainsAll([]string{"xxx", "zzz"})
	if found {
		t.Fatal("Unexpected")
	}
}
