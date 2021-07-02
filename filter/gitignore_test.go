package filter_test

import (
	"testing"

	"github.com/komem3/fing/filter"
)

func TestGitignore_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name     string
		ignore   []*filter.Path
		filename string
		match    bool
	}{
		{
			"match",
			[]*filter.Path{
				filter.NewPath("coverage.*"),
			},
			"coverage.out",
			true,
		},
		{
			"negative case",
			[]*filter.Path{
				filter.NewPath("node_modules/*"),
				filter.NewNotPath("*index.js"),
			},
			"node_modules/sample/index.js",
			false,
		},
		{
			"negative negative case",
			[]*filter.Path{
				filter.NewPath("node_modules/*"),
				filter.NewNotPath("*index.js"),
				filter.NewPath("node_modules/sample/*"),
			},
			"node_modules/sample/index.js",
			true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if match, _ := (&filter.Gitignore{PathMatchers: tt.ignore}).Match(tt.filename, nil); match != tt.match {
				t.Errorf("Match want %t, but got %t", tt.match, match)
			}
		})
	}
}
