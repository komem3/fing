package filter

import (
	"fmt"
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

func (f *FileName) Match(_ string, info fs.DirEntry) bool {
	return f.match(info.Name())
}

func (f IFileName) Match(_ string, info fs.DirEntry) bool {
	return f.match(strings.ToUpper(info.Name()))
}

func (f *FileName) String() string {
	return fmt.Sprintf("name(%s)", f.glob)
}

func (f *IFileName) String() string {
	return fmt.Sprintf("iname(%s)", f.glob)
}
