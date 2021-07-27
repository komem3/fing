package filter_test

import (
	"io/fs"
	"reflect"
	"testing"

	"github.com/komem3/fing/filter"
)

func TestNewSize(t *testing.T) {
	for _, tt := range []struct {
		arg    string
		size   *filter.Size
		errMsg string
	}{
		{"+3c", &filter.Size{3, filter.GreaterCmpOption}, ""},
		{"-3c", &filter.Size{3, filter.LessCmpOption}, ""},
		{"1k", &filter.Size{1024, filter.EqualCmpOption}, ""},
		{"1M", &filter.Size{1048576, filter.EqualCmpOption}, ""},
		{"1G", &filter.Size{1073741824, filter.EqualCmpOption}, ""},
		{"1", nil, "1 is invalid unit of size"},
		{"1m", nil, "m is invalid unit of size"},
		{"+", nil, "+ is invalid size argument"},
		{"", nil, "missing argument of size"},
	} {
		tt := tt
		t.Run(tt.arg, func(t *testing.T) {
			t.Parallel()
			size, err := filter.NewSize(tt.arg)
			if want, got := (tt.errMsg != ""), err != nil; want != got {
				t.Errorf("err != nil want %t, but got %t", want, got)
			}
			if tt.errMsg != "" && tt.errMsg != err.Error() {
				t.Errorf("err.Error() mismatch\nwant: %s\ngot: %s", tt.errMsg, err.Error())
			}
			if !reflect.DeepEqual(tt.size, size) {
				t.Errorf("NewSize() mismatch\nwant: %#v\ngot: %#v", tt.size, size)
			}
		})
	}
}

func TestSize_Match(t *testing.T) {
	for _, tt := range []struct {
		name  string
		size  *filter.Size
		arg   fs.DirEntry
		match bool
	}{
		{
			"equal", &filter.Size{1024, filter.EqualCmpOption},
			&mockDirFileInfo{size: 1024}, true,
		},
		{
			"greater", &filter.Size{1024, filter.GreaterCmpOption},
			&mockDirFileInfo{size: 1025}, true,
		},
		{
			"less", &filter.Size{1024, filter.LessCmpOption},
			&mockDirFileInfo{size: 1023}, true,
		},
		{
			"not greater", &filter.Size{1024, filter.GreaterCmpOption},
			&mockDirFileInfo{size: 1024}, false,
		},
		{
			"not less", &filter.Size{1024, filter.LessCmpOption},
			&mockDirFileInfo{size: 1024}, false,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if match, _ := tt.size.Match("", tt.arg); tt.match != match {
				t.Errorf("Match() mistamch want %t, but got %t", tt.match, match)
			}
		})
	}
}
