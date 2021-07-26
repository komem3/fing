package filter_test

import (
	"testing"

	"github.com/komem3/fing/filter"
)

const pathFile = "src/index.ts"

func TestPath_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern string
		match   bool
	}{
		{"src/*.ts", true},
		{"src/*.TS", false},
		{"index.ts", false},
	} {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			filter := filter.NewPath(tt.pattern)
			if match, _ := filter.Match(pathFile, nil); match != tt.match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}

func TestIPath_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern string
		match   bool
	}{
		{"src/*.ts", true},
		{"src/*.TS", true},
		{"index.ts", false},
	} {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			filter := filter.NewIPath(tt.pattern)
			if match, _ := filter.Match(pathFile, nil); match != tt.match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}
