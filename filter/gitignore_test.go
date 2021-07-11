package filter_test

import (
	"reflect"
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

func TestGitignore_Add(t *testing.T) {
	for _, tt := range []struct {
		name string
		dist *filter.Gitignore
		src  *filter.Gitignore
		want *filter.Gitignore
	}{
		{
			"merge gitignore",
			&filter.Gitignore{[]*filter.Path{filter.NewPath("*.txt")}},
			&filter.Gitignore{[]*filter.Path{filter.NewPath("*.png")}},
			&filter.Gitignore{[]*filter.Path{filter.NewPath("*.txt"), filter.NewPath("*.png")}},
		},
		{
			"dist is empty",
			nil,
			&filter.Gitignore{[]*filter.Path{filter.NewPath("*.png")}},
			&filter.Gitignore{[]*filter.Path{filter.NewPath("*.png")}},
		},
		{
			"src is empty",
			&filter.Gitignore{[]*filter.Path{filter.NewPath("*.txt")}},
			nil,
			&filter.Gitignore{[]*filter.Path{filter.NewPath("*.txt")}},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			distLen, srcLen := tt.dist.Len(), tt.src.Len()
			got := tt.dist.Add(tt.src)
			if distLen != tt.dist.Len() {
				t.Errorf("change dist length %d -> %d", distLen, len(tt.dist.PathMatchers))
			}
			if srcLen != tt.src.Len() {
				t.Errorf("change src length %d -> %d", srcLen, len(tt.src.PathMatchers))
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Add want -, got +\n-%s\n+%s", tt.want, got)
			}
		})
	}
}
