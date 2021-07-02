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
			args{filter.NewMockDriEntry("mock.go", false, 0, nil), "f"},
			want{true, false},
		},
		{
			"match directory",
			args{filter.NewMockDriEntry("mock/", true, fs.ModeDir, nil), "d"},
			want{true, false},
		},
		{
			"match pipe",
			args{filter.NewMockDriEntry("pipe", false, fs.ModeNamedPipe, nil), "p"},
			want{true, false},
		},
		{
			"match socket",
			args{filter.NewMockDriEntry("socket", false, fs.ModeSocket, nil), "s"},
			want{true, false},
		},
		{
			"mismatch directory",
			args{filter.NewMockDriEntry("mock/", true, fs.ModeDir, nil), "f"},
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
