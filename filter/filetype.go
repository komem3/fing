package filter

import (
	"fmt"
	"io/fs"
)

type FileType fs.FileMode

var _ FileExp = FileType(0)

func NewFileType(typ string) (FileType, error) {
	switch typ {
	case "f":
		return 0, nil
	case "d":
		return FileType(fs.ModeDir), nil
	case "p":
		return FileType(fs.ModeNamedPipe), nil
	case "s":
		return FileType(fs.ModeSocket), nil
	}
	return 0, fmt.Errorf("%s is invalid file type", typ)
}

func (f FileType) Match(_ string, info fs.DirEntry) bool {
	return fs.FileMode(f) == info.Type()
}

func (f FileType) String() string {
	if fs.FileMode(f).IsRegular() {
		return "type(file)"
	}
	return fmt.Sprintf("type(%s)", fs.FileMode(f))
}
