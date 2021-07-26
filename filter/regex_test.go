package filter_test

import (
	"testing"

	"github.com/komem3/fing/filter"
)

const regexFile = "src/index.ts"

func TestRegex_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern string
		match   bool
	}{
		{`src/.*\.ts`, true},
		{`src/.*\.TS`, false},
		{`src/false.ts`, false},
	} {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			filter, err := filter.NewRegex(tt.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if match, _ := filter.Match(regexFile, nil); match != tt.match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}

func TestIRegex_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern string
		match   bool
	}{
		{`src/.*\.ts`, true},
		{`src/.*\.TS`, true},
		{`src/false.ts`, false},
	} {
		tt := tt
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			filter, err := filter.NewIRegex(tt.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if match, _ := filter.Match(regexFile, nil); match != tt.match {
				t.Errorf("match want %t, but got %t", tt.match, match)
			}
		})
	}
}
