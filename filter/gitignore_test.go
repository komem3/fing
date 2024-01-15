package filter_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/komem3/fing/filter"
)

var tmpDir = os.TempDir()

var domains = strings.Split(tmpDir, string(filepath.Separator))

var extactTests = []struct {
	name   string
	data   []string
	ignore *filter.Gitignore
}{
	{
		"parse gitignore",
		[]string{
			"# comment out",
			"node_modules/**",
			"/vendor",
			"*.jpg",
			"!*.txt",
			"!sample.png",
			"!/root.png",
		},
		&filter.Gitignore{
			PathMatchers: []gitignore.Pattern{
				gitignore.ParsePattern(filepath.ToSlash("node_modules/**"), domains),
				gitignore.ParsePattern(filepath.ToSlash("/vendor"), domains),
				gitignore.ParsePattern("*.jpg", domains),
				gitignore.ParsePattern("!*.txt", domains),
				gitignore.ParsePattern("!sample.png", domains),
				gitignore.ParsePattern("!"+filepath.ToSlash("/root.png"), domains),
			},
		},
	},
	{
		"empty file",
		[]string{},
		&filter.Gitignore{},
	},
}

func TestNewGitIgnore(t *testing.T) {
	for _, tt := range extactTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tmp, err := os.CreateTemp(tmpDir, "testgitignore*")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { os.Remove(tmp.Name()) })
			if _, err := tmp.WriteString(strings.Join(tt.data, "\n")); err != nil {
				t.Fatal(err)
			}
			tmp.Close()

			ignores, err := filter.NewGitIgnore(tmpDir, filepath.Join(tmpDir, filepath.Base(tmp.Name())))
			if err != nil {
				t.Fatal(err)
			}
			if len(tt.ignore.PathMatchers) != len(ignores.PathMatchers) {
				t.Errorf("ignores len mismatch want -, but got +\n-%s\n+%s", tt.ignore.PathMatchers, ignores.PathMatchers)
			}
			for i, ignore := range tt.ignore.PathMatchers {
				if !reflect.DeepEqual(ignore, ignores.PathMatchers[i]) {
					t.Errorf("ignore mismatch want -, but got +\n-%s\n+%s", ignore, ignores.PathMatchers[i])
				}
			}
		})
	}
}

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
				gitignore.ParsePattern("/*.out", []string{"root", "test"}),
			},
			filepath.FromSlash("root/test/coverage.out"),
			true,
		},
		{
			"negative case",
			[]gitignore.Pattern{
				gitignore.ParsePattern("node_modules/**", []string{"root"}),
				gitignore.ParsePattern("!**/index.js", []string{"root"}),
			},
			filepath.FromSlash("root/node_modules/sample/index.js"),
			false,
		},
		{
			"negative negative case",
			[]gitignore.Pattern{
				gitignore.ParsePattern("node_modules/**", []string{"root"}),
				gitignore.ParsePattern("!**/index.js", []string{"root"}),
				gitignore.ParsePattern("node_modules/sample/*", []string{"root"}),
			},
			filepath.FromSlash("root/node_modules/sample/index.js"),
			true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if match, _ := (&filter.Gitignore{PathMatchers: tt.ignore}).
				Match(tt.filename, &mockDirFileInfo{isDir: false}); match != tt.match {
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
