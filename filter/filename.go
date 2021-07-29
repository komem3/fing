package filter

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type (
	FileName  string
	IFileName string
)

var (
	_ FileExp = FileName("")
	_ FileExp = IFileName("")
)

func NewFileName(pattern string) FileName {
	return FileName(pattern)
}

func NewIFileName(pattern string) IFileName {
	return IFileName(strings.ToUpper(pattern))
}

func (f FileName) Match(_ string, info fs.DirEntry) (bool, error) {
	return filepath.Match(string(f), info.Name())
}

func (f IFileName) Match(_ string, info fs.DirEntry) (bool, error) {
	return filepath.Match(string(f), strings.ToUpper(info.Name()))
}

func (f FileName) String() string {
	return fmt.Sprintf("name(%s)", string(f))
}

func (f IFileName) String() string {
	return fmt.Sprintf("iname(%s)", string(f))
}
