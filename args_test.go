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
	t.Log(args)
	if reflect.DeepEqual(args.words, src) {
		t.Error("words links to the soruce slice")
	}
}

func TestArgsCurrent(t *testing.T) {
	for expect, args := range map[string]Args{
		"zero": Args{
			words: []string{"zero", "one", "two"},
			index: 0,
		},
		"one": Args{
			words: []string{"zero", "one", "two"},
			index: 1,
		},
		"two": Args{
			words: []string{"zero", "one", "two"},
			index: 2,
		},
	} {
		t.Run(expect, func(t *testing.T) {
			actual := args.Current()
			if expect != actual {
				t.Errorf("current arg is not match. expect: %v, actual: %v", expect, actual)
			}
		})
	}
}

func TestArgsNext(t *testing.T) {
	for i, args := range map[int]Args{
		1: Args{
			words: []string{"zero", "one", "two"},
			index: 0,
		},
		2: Args{
			words: []string{"zero", "one", "two"},
			index: 1,
		},
		-1: Args{
			words: []string{"zero", "one", "two"},
			index: 2,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			next := args.Next()
			if i < 0 {
				if next != nil {
					t.Errorf("next should be nil, but actual %v", next)
				}
			} else {
				if i != next.index {
					t.Errorf("the index of the next args is mismatch. expect: %v, actual: %v", i, next.index)
				}
			}
		})
	}
}
