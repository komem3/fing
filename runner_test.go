package main

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var tests = []struct {
	command string
	output  []string
}{
	{
		"fing testdata/jpg_dir testdata/png_dir",
		[]string{
			"testdata/jpg_dir/1.jpg",
			"testdata/jpg_dir/2.jpg",
			"testdata/jpg_dir/3.jpg",
			"testdata/jpg_dir/4.JPG",
			"testdata/png_dir/1.png",
			"testdata/png_dir/2.png",
			"testdata/png_dir/3.png",
		},
	},
	{
		"fing -I testdata -type f",
		[]string{
			"testdata/.gitignore",
			"testdata/txt_dir/1.txt",
			"testdata/txt_dir/2.txt",
			"testdata/txt_dir/.gitignore",
		},
	},
	{
		"fing testdata -iname *.jpg -regex (3|4).*",
		[]string{
			"testdata/jpg_dir/3.jpg",
			"testdata/jpg_dir/4.JPG",
		},
	},
	{
		"fing testdata -path *_dir -prune -iregex (1|2).*",
		[]string{
			"testdata/link/1.ln",
			"testdata/link/2.ln",
		},
	},
	{
		"fing testdata -ipath *_dir/* -not -name *.jpg",
		[]string{
			"testdata/jpg_dir/4.JPG",
			"testdata/png_dir/1.png",
			"testdata/png_dir/2.png",
			"testdata/png_dir/3.png",
		},
	},
}

func TestRun(t *testing.T) {
	for _, tt := range tests {
		tt := tt
		t.Run(tt.command, func(t *testing.T) {
			t.Parallel()
			out := new(bytes.Buffer)
			outerr := new(bytes.Buffer)
			if err := run(strings.Split(tt.command, " "), out, outerr); err != nil {
				t.Errorf("error: %s", outerr.String())
			}
			if errStr := outerr.String(); errStr != "" {
				t.Error(errStr)
			}

			files := strings.Split(out.String(), "\n")
			if files[len(files)-1] == "" {
				files = files[:len(files)-1]
			}
			sort.Strings(files)
			sort.Strings(tt.output)
			if !reflect.DeepEqual(tt.output, files) {
				t.Errorf("output is mismatch\nwant:\n%s\ngot:\n%s",
					strings.Join(tt.output, "\n"),
					strings.Join(files, "\n"),
				)
			}
		})
	}
}
