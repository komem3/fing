package walk

import (
	"io/fs"
	"os"
)

func newEntry(f *os.File) (fs.DirEntry, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return fs.FileInfoToDirEntry(info), nil
}
