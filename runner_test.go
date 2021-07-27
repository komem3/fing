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
		"fing testdata/jpg_dir testdata/png_dir -empty",
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
		"fing testdata -I -type f -not -name .*",
		[]string{
			"testdata/txt_dir/1.txt",
			"testdata/txt_dir/2.txt",
		},
	},
	{
		"fing testdata -iname *.jpg -regex .*(3|4).*",
		[]string{
			"testdata/jpg_dir/3.jpg",
			"testdata/jpg_dir/4.JPG",
		},
	},
	{
		"fing testdata -name .* -prune -path */link/* -or -name *.txt",
		[]string{
			"testdata/link/1.ln",
			"testdata/link/2.ln",
			"testdata/txt_dir/1.txt",
			"testdata/txt_dir/2.txt",
		},
	},
	{
		"fing testdata -not -name .* -prune -not -type f",
		[]string{
			"testdata/.hidden",
		},
	},
	{
		"fing testdata -name jpg* -or -name png* -prune -irname (1|2).*",
		[]string{
			"testdata/link/1.ln",
			"testdata/link/2.ln",
			"testdata/txt_dir/1.txt",
			"testdata/txt_dir/2.txt",
		},
	},
	{
		"fing testdata -type f -ipath txt",
		[]string{},
	},
	{
		"fing testdata -maxdepth 1 -name .gitignore",
		[]string{
			"testdata/.gitignore",
		},
	},
	{
		"fing testdata -maxdepth 0",
		[]string{
			"testdata",
		},
	},
	{
		"fing testdata/jpg_dir testdata/png_dir -dry -I -type f -ipath txt/* -prune -name *.png -or -not -regex .*\\.name",
		[]string{
			"targets=[testdata/jpg_dir, testdata/png_dir] " +
				"ignore=true prunes=[type(file) * ipath(TXT/*)] " +
				"condition=[name(*.png) + not regex(^.*\\.name$)]",
		},
	},
	{
		"fing testdata -size +1k",
		[]string{
			"testdata",
			"testdata/txt_dir",
			"testdata/.hidden",
			"testdata/link",
			"testdata/png_dir",
			"testdata/jpg_dir",
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
