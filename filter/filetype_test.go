package filter_test

import (
	"io/fs"
	"testing"

	"github.com/komem3/fing/filter"
)

func TestFileType_Match(t *testing.T) {
	type (
		args struct {
			file fs.DirEntry
			typ  string
		}
		want struct {
			match bool
			isErr bool
		}
	)
	for _, tt := range []struct {
		name string
		args args
		want want
	}{
		{
			"invalid type",
			args{nil, "invalid"},
			want{false, true},
		},
		{
			"match regular file",
			args{&mockDirFileInfo{typ: 0}, "f"},
			want{true, false},
		},
		{
			"match directory",
			args{&mockDirFileInfo{typ: fs.ModeDir}, "d"},
			want{true, false},
		},
		{
			"match pipe",
			args{&mockDirFileInfo{typ: fs.ModeNamedPipe}, "p"},
			want{true, false},
		},
		{
			"match socket",
			args{&mockDirFileInfo{typ: fs.ModeSocket}, "s"},
			want{true, false},
		},
		{
			"mismatch directory",
			args{&mockDirFileInfo{typ: fs.ModeDir}, "f"},
			want{false, false},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f, err := filter.NewFileType(tt.args.typ)
			if got := err != nil; got != tt.want.isErr {
				t.Errorf("err != nil want %t, but got \n%v", tt.want.isErr, err)
			}
			if err != nil {
				return
			}
			if match, _ := f.Match("", tt.args.file); match != tt.want.match {
				t.Errorf("match want %t, but got %t", tt.want.match, match)
			}
		})
	}
}
