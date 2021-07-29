package walk

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/komem3/fing/filter"
)

var tmpDir = os.TempDir()

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
				gitignore.ParsePattern(filepath.ToSlash(filepath.Join(tmpDir, "node_modules/**")), nil),
				gitignore.ParsePattern(filepath.ToSlash(filepath.Join(tmpDir, "vendor")), nil),
				gitignore.ParsePattern("*.jpg", nil),
				gitignore.ParsePattern("!*.txt", nil),
				gitignore.ParsePattern("!sample.png", nil),
				gitignore.ParsePattern("!"+filepath.ToSlash(filepath.Join(tmpDir, "root.png")), nil),
			},
		},
	},
	{
		"empty file",
		[]string{},
		&filter.Gitignore{
			PathMatchers: make([]gitignore.Pattern, 0, defaultIgnoreBuffer),
		},
	},
}

func TestWalker_extractGitignore(t *testing.T) {
	walker := &Walker{}
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

			ignore, err := walker.extractGitignore(tmpDir, filepath.Base(tmp.Name()))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.ignore, ignore) {
				t.Errorf("extractGitignore want -, but got +\n-%s\n+%s", tt.ignore, ignore)
			}
		})
	}
}

var testFs = fstest.MapFS{
	"testdata/.gitignore":         {},
	"testdata/sample.txt":         {},
	"testdata/example.png":        {},
	"testdata/jpg_dir/.gitignore": {},
	"testdata/jpg_dir/sample.jpg": {},
	"testdata/link/sample.ln":     {},
}

func TestWalker_getIgnore(t *testing.T) {
	walker := &Walker{}
	for _, tt := range []struct {
		name   string
		path   string
		result string
	}{
		{"get gitignore", "testdata", ".gitignore"},
		{"not contain", "testdata/link", ""},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			entries, err := testFs.ReadDir(tt.path)
			if err != nil {
				t.Fatal(err)
			}
			ignore := walker.getIgnore(entries)
			if ignore != tt.result {
				t.Errorf("getIgnore want %s, but got %s", tt.result, ignore)
			}
		})
	}
}
