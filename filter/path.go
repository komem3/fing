package filter

import (
	"fmt"
	"io/fs"
	"strings"
)

type (
	Path struct {
		*glob
		pathType pathType

		fmt.Stringer
	}
	IPath struct {
		*glob
		pathType pathType
	}
)

type pathType int

const (
	mismatchPathType pathType = iota - 1
	normalPathType
	notPathType
)

var (
	_ FileExp = (*Path)(nil)
	_ FileExp = (*IPath)(nil)
)

func NewPath(pattern string) *Path {
	return &Path{
		glob:     newGlob(pattern),
		pathType: normalPathType,
	}
}

func NewNotPath(pattern string) *Path {
	return &Path{
		glob:     newGlob(pattern),
		pathType: notPathType,
	}
}

func NewIPath(pattern string) *IPath {
	return &IPath{
		glob: newGlob(strings.ToUpper(pattern)),
	}
}

func (p *Path) Match(path string, _ fs.DirEntry) (bool, error) {
	return p.match(path), nil
}

func (p *IPath) Match(path string, _ fs.DirEntry) (bool, error) {
	return p.match(strings.ToUpper(path)), nil
}

func (p *Path) String() string {
	return fmt.Sprintf("%s(%d)", p.glob, p.pathType)
}
