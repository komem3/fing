package filter_test

import (
	"testing"

	"github.com/komem3/fing/filter"
)

var regexNameFile = &mockDirFileInfo{name: "index.ts"}

func TestRegexName_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern string
		match   bool
	}{
		{".*\\.ts", true},
		{".*\\.TS", false},
		{"src/.*ts", false},
	} {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			filter, err := filter.NewRegexName(tt.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if match, _ := filter.Match("", regexNameFile); match != tt.match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}

func TestIRegexName_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern string
		match   bool
	}{
		{".*\\.ts", true},
		{".*\\.TS", true},
		{"src/.*ts", false},
	} {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			filter, err := filter.NewIRegexName(tt.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if match, _ := filter.Match("", regexNameFile); match != tt.match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}
