package walk

import (
	"os"
	"testing"
	"testing/fstest"
)

var tmpDir = os.TempDir()

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
