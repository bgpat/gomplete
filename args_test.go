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
		t.Errorf("words is not match to src. want: %v, got: %v", src, args.words)
	}
	src[1] = "awesome"
	if reflect.DeepEqual(args.words, src) {
		t.Error("words links to the soruce slice")
	}
}

func TestArgsCurrent(t *testing.T) {
	for want, args := range map[string]Args{
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
		t.Run(want, func(t *testing.T) {
			got := args.Current()
			if want != got {
				t.Errorf("Current argument is invalid. want: %v, got: %v", want, got)
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
					t.Errorf("The next must be nil, but got %v.", next)
				}
			} else if next == nil {
				t.Errorf("The next must not be nil. want: %v", i)
			} else {
				if i != next.index {
					t.Errorf("The index of the next args is mismatch. want: %v, got: %v", i, next.index)
				}
			}
		})
	}
}
