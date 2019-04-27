package gomplete

import (
	"reflect"
	"strconv"
	"testing"
)

func TestNewArgs(t *testing.T) {
	src := []string{"foo", "bar", "baz"}
	args := NewArgs(src)
	if !reflect.DeepEqual(args.words, src) {
		t.Errorf("words is not match to src. expect: %v, actual: %v", src, args.words)
	}
	src[1] = "awesome"
	if reflect.DeepEqual(args.words, src) {
		t.Error("words links to the soruce slice")
	}
}

func TestArgsCurrent(t *testing.T) {
	for expect, args := range map[string]Args{
		"zero": {
			words: []string{"zero", "one", "two"},
			index: 0,
		},
		"one": {
			words: []string{"zero", "one", "two"},
			index: 1,
		},
		"two": {
			words: []string{"zero", "one", "two"},
			index: 2,
		},
	} {
		t.Run(expect, func(t *testing.T) {
			actual := args.Current()
			if expect != actual {
				t.Errorf("Current argument is invalid. expect: %v, actual: %v", expect, actual)
			}
		})
	}
}

func TestArgsNext(t *testing.T) {
	for i, args := range map[int]Args{
		1: {
			words: []string{"zero", "one", "two"},
			index: 0,
		},
		2: {
			words: []string{"zero", "one", "two"},
			index: 1,
		},
		-1: {
			words: []string{"zero", "one", "two"},
			index: 2,
		},
		-2: {
			words: []string{},
			index: 0,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			next := args.Next()
			if i < 0 {
				if next != nil {
					t.Errorf("The next must be nil, but actual %v.", next)
				}
			} else if next == nil {
				t.Errorf("The next must not be nil. expect: %v", i)
			} else {
				if i != next.index {
					t.Errorf("The index of the next args is mismatch. expect: %v, actual: %v", i, next.index)
				}
			}
		})
	}
}
