package main

import (
	"bytes"
	"path/filepath"
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
		"fing testdata/jpg_dir/1.jpg",
		[]string{
			"testdata/jpg_dir/1.jpg",
		},
	},
	{
		"fing testdata/jpg_dir testdata/png_dir -empty -type f",
		[]string{
			filepath.FromSlash("testdata/jpg_dir/1.jpg"),
			filepath.FromSlash("testdata/jpg_dir/2.jpg"),
			filepath.FromSlash("testdata/jpg_dir/3.jpg"),
			filepath.FromSlash("testdata/jpg_dir/4.JPG"),
			filepath.FromSlash("testdata/png_dir/1.png"),
			filepath.FromSlash("testdata/png_dir/2.png"),
			filepath.FromSlash("testdata/png_dir/3.png"),
		},
	},
	{
		"fing testdata -I -type f -not -name .*",
		[]string{
			filepath.FromSlash("testdata/scripts/test.sh"),
			filepath.FromSlash("testdata/scripts/README.md"),
			filepath.FromSlash("testdata/txt_dir/1.txt"),
			filepath.FromSlash("testdata/txt_dir/2.txt"),
		},
	},
	{
		"fing testdata -iname *.jpg -regex .*(3|4).*",
		[]string{
			filepath.FromSlash("testdata/jpg_dir/3.jpg"),
			filepath.FromSlash("testdata/jpg_dir/4.JPG"),
		},
	},
	{
		"fing testdata -name .* -prune -false -o -path */link/* -o -name *.txt",
		[]string{
			filepath.FromSlash("testdata/link/1.ln"),
			filepath.FromSlash("testdata/link/2.ln"),
			filepath.FromSlash("testdata/txt_dir/1.txt"),
			filepath.FromSlash("testdata/txt_dir/2.txt"),
		},
	},
	{
		"fing testdata -rname ^(jpg|png) -prune -false -o -irname (1|2).*",
		[]string{
			filepath.FromSlash("testdata/link/1.ln"),
			filepath.FromSlash("testdata/link/2.ln"),
			filepath.FromSlash("testdata/txt_dir/1.txt"),
			filepath.FromSlash("testdata/txt_dir/2.txt"),
		},
	},
	{
		"fing testdata -type f -ipath txt",
		[]string{},
	},
	{
		"fing -maxdepth 1 -name .gitignore testdata",
		[]string{
			filepath.FromSlash("testdata/.gitignore"),
		},
	},
	{
		"fing testdata -maxdepth 0",
		[]string{
			"testdata",
		},
	},
	{
		"fing testdata/ -executable -type d -name scripts",
		[]string{
			filepath.FromSlash("testdata/scripts"),
		},
	},
	{
		"fing testdata/jpg_dir testdata/png_dir -dry -I -iname txt* -prune -false -o -name *.png -o -not -regex .*\\.name",
		[]string{
			"targets=[testdata/jpg_dir, testdata/png_dir] " +
				filepath.FromSlash("ignore=true prunes=[iname(TXT*)] ") +
				"condition=[iname(TXT*) && false || name(*.png) || not regex(.*\\.name)]",
		},
	},
	{
		"fing testdata -size +0c -type f",
		[]string{
			filepath.FromSlash("testdata/.gitignore"),
			filepath.FromSlash("testdata/scripts/test.sh"),
			filepath.FromSlash("testdata/scripts/README.md"),
			filepath.FromSlash("testdata/txt_dir/.gitignore"),
			filepath.FromSlash("testdata/txt_dir/1.txt"),
			filepath.FromSlash("testdata/txt_dir/2.txt"),
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
			if status := run(strings.Split(tt.command, " "), out, outerr); status != 0 {
				t.Fatalf("error status: %s", outerr)
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
