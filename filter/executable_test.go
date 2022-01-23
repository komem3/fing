package filter_test

import (
	"io/fs"
	"testing"

	"github.com/komem3/fing/filter"
)

func TestExecutable_Match(t *testing.T) {
	for _, tt := range []struct {
		name  string
		info  fs.DirEntry
		match bool
	}{
		{"executable", &mockDirFileInfo{typ: 100}, true},
		{"not executable", &mockDirFileInfo{typ: 11}, false},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if match, _ := (&filter.Executable{}).Match("", tt.info); tt.match != match {
				t.Errorf("Executable.Match mismatch want %t, but got %t", tt.match, match)
			}
		})
	}
}
