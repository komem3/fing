package filter_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/komem3/fing/filter"
)

func TestGitignore_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name     string
		ignore   []gitignore.Pattern
		filename string
		match    bool
	}{
		{
			"match",
			[]gitignore.Pattern{
				gitignore.ParsePattern("/coverage.*", nil),
			},
			"coverage.out",
			true,
		},
		{
			"negative case",
			[]gitignore.Pattern{
				gitignore.ParsePattern("node_modules/**", nil),
				gitignore.ParsePattern("!**/index.js", nil),
			},
			filepath.FromSlash("node_modules/sample/index.js"),
			false,
		},
		{
			"negative negative case",
			[]gitignore.Pattern{
				gitignore.ParsePattern("node_modules/**", nil),
				gitignore.ParsePattern("!**/index.js", nil),
				gitignore.ParsePattern("node_modules/sample/*", nil),
			},
			filepath.FromSlash("node_modules/sample/index.js"),
			true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if match, _ := (&filter.Gitignore{PathMatchers: tt.ignore}).Match(tt.filename, &mockDirFileInfo{isDir: false}); match != tt.match {
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
			&filter.Gitignore{[]gitignore.Pattern{gitignore.ParsePattern("*.txt", nil)}},
			&filter.Gitignore{[]gitignore.Pattern{gitignore.ParsePattern("*.png", nil)}},
			&filter.Gitignore{[]gitignore.Pattern{gitignore.ParsePattern("*.txt", nil), gitignore.ParsePattern("*.png", nil)}},
		},
		{
			"dist is empty",
			nil,
			&filter.Gitignore{[]gitignore.Pattern{gitignore.ParsePattern("*.png", nil)}},
			&filter.Gitignore{[]gitignore.Pattern{gitignore.ParsePattern("*.png", nil)}},
		},
		{
			"src is empty",
			&filter.Gitignore{[]gitignore.Pattern{gitignore.ParsePattern("*.txt", nil)}},
			nil,
			&filter.Gitignore{[]gitignore.Pattern{gitignore.ParsePattern("*.txt", nil)}},
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
