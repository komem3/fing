package filter_test

import (
	"testing"

	"github.com/komem3/fing/filter"
)

func TestOrExp_Match(t *testing.T) {
	const path = "test.txt"
	for _, tt := range []struct {
		name   string
		filter filter.OrExp
		match  bool
	}{
		{
			"match partial",
			filter.OrExp{filter.NewPath("test.*"), filter.NewPath("*.jpg")},
			true,
		},
		{
			"mismatch all",
			filter.OrExp{filter.NewPath("*.png"), filter.NewPath("*.jpg")},
			false,
		},
		{
			"empty expression",
			filter.OrExp{},
			true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			match, err := tt.filter.Match(path, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.match != match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}

func TestAndExp_Match(t *testing.T) {
	const path = "test.txt"
	for _, tt := range []struct {
		name   string
		filter filter.AndExp
		match  bool
	}{
		{
			"match all",
			filter.AndExp{filter.NewPath("test.*"), filter.NewPath("*.txt")},
			true,
		},
		{
			"match partial",
			filter.AndExp{filter.NewPath("test.*"), filter.NewPath("*.jpg")},
			false,
		},
		{
			"empty expression",
			filter.AndExp{},
			true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			match, err := tt.filter.Match(path, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.match != match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}

func TestNotFilter_Match(t *testing.T) {
	for _, tt := range []struct {
		name   string
		path   string
		filter *filter.NotExp
		match  bool
	}{
		{
			"match path",
			"test.txt",
			filter.NewNotExp(filter.NewPath("test.txt")),
			false,
		},
		{
			"mismatch path",
			"test.txt",
			filter.NewNotExp(filter.NewPath("miss.txt")),
			true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			match, err := tt.filter.Match(tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.match != match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}
