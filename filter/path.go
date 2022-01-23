package filter

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/komem3/glob"
)

type (
	Path  struct{ *glob.Glob }
	IPath struct{ *glob.Glob }
)

var (
	_ FileExp = (*Path)(nil)
	_ FileExp = (*IPath)(nil)
)

func NewPath(pattern string) (*Path, error) {
	glob, err := glob.Compile(escapeBackSlash(pattern))
	if err != nil {
		return nil, err
	}
	return &Path{Glob: glob}, nil
}

func NewIPath(pattern string) (*IPath, error) {
	glob, err := glob.Compile(strings.ToUpper(escapeBackSlash(pattern)))
	if err != nil {
		return nil, err
	}
	return &IPath{Glob: glob}, nil
}

func (p Path) Match(path string, _ fs.DirEntry) (bool, error) {
	return p.MatchString(path), nil
}

func (p IPath) Match(path string, _ fs.DirEntry) (bool, error) {
	return p.MatchString(strings.ToUpper(path)), nil
}

func (p Path) String() string {
	return fmt.Sprintf("path(%s)", p.Glob)
}

func (p IPath) String() string {
	return fmt.Sprintf("ipath(%s)", p.Glob)
}
