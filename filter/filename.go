package filter

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/komem3/glob"
)

type (
	FileName  struct{ *glob.Glob }
	IFileName struct{ *glob.Glob }
)

var (
	_ FileExp = (*FileName)(nil)
	_ FileExp = (*FileName)(nil)
)

func NewFileName(pattern string) (*FileName, error) {
	glob, err := glob.Compile(escapeBackSlash(pattern))
	if err != nil {
		return nil, err
	}
	return &FileName{Glob: glob}, nil
}

func NewIFileName(pattern string) (*IFileName, error) {
	glob, err := glob.Compile(strings.ToUpper(escapeBackSlash(pattern)))
	if err != nil {
		return nil, err
	}
	return &IFileName{Glob: glob}, nil
}

func (f FileName) Match(_ string, info fs.DirEntry) (bool, error) {
	return f.MatchString(info.Name()), nil
}

func (f IFileName) Match(_ string, info fs.DirEntry) (bool, error) {
	return f.MatchString(strings.ToUpper(info.Name())), nil
}

func (f FileName) String() string {
	return fmt.Sprintf("name(%s)", f.Glob)
}

func (f IFileName) String() string {
	return fmt.Sprintf("iname(%s)", f.Glob)
}
