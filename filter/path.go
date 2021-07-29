package filter

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type (
	Path  string
	IPath string
)

var (
	_ FileExp = Path("")
	_ FileExp = IPath("")
)

func NewPath(pattern string) Path {
	return Path(pattern)
}

func NewIPath(pattern string) IPath {
	return IPath(strings.ToUpper(pattern))
}

func (p Path) Match(path string, _ fs.DirEntry) (bool, error) {
	return filepath.Match(string(p), path)
}

func (p IPath) Match(path string, _ fs.DirEntry) (bool, error) {
	return filepath.Match(string(p), strings.ToUpper(path))
}

func (p Path) String() string {
	return fmt.Sprintf("path(%s)", string(p))
}

func (p IPath) String() string {
	return fmt.Sprintf("ipath(%s)", string(p))
}
