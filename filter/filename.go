package filter

import (
	"io/fs"
	"strings"
)

type (
	FileName struct {
		*glob
	}
	IFileName struct {
		*glob
	}
)

var (
	_ FileExp = (*FileName)(nil)
	_ FileExp = (*IFileName)(nil)
)

func NewFileName(pattern string) *FileName {
	return &FileName{glob: newGlob(pattern)}
}

func NewIFileName(pattern string) *IFileName {
	return &IFileName{glob: newGlob(strings.ToUpper(pattern))}
}

func (f *FileName) Match(_ string, info fs.DirEntry) (bool, error) {
	return f.match(info.Name()), nil
}

func (f IFileName) Match(_ string, info fs.DirEntry) (bool, error) {
	return f.match(strings.ToUpper(info.Name())), nil
}
