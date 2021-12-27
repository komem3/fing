package filter_test

import (
	"testing"

	"github.com/komem3/fing/filter"
)

var fileNameFile = &mockDirFileInfo{name: "index.ts"}

func TestFileName_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern string
		match   bool
	}{
		{"index.ts", true},
		{"index.TS", false},
		{"src/*.ts", false},
	} {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			filter, err := filter.NewFileName(tt.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if match, _ := filter.Match("", fileNameFile); match != tt.match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}

func TestIFileName_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern string
		match   bool
	}{
		{"index.ts", true},
		{"index.TS", true},
		{"src/*.ts", false},
	} {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			filter, err := filter.NewIFileName(tt.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if match, _ := filter.Match("", fileNameFile); match != tt.match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}
