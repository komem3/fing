package walk

import (
	"io/fs"
	"os"
)

type dirEntry struct {
	fs.FileInfo
}

func newEntry(f *os.File) (*dirEntry, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return &dirEntry{info}, nil
}

var _ fs.DirEntry = (*dirEntry)(nil)

func (d *dirEntry) Info() (fs.FileInfo, error) {
	return d.FileInfo, nil
}

func (d *dirEntry) Type() fs.FileMode {
	return d.Mode()
}
