package filter

import (
	"io/fs"
)

type Executable struct{}

var _ FileExp = (*Executable)(nil)

const executablePerm = 0o100

func NewExecutable() *Executable {
	return new(Executable)
}

func (*Executable) Match(_ string, info fs.DirEntry) (bool, error) {
	inf, err := info.Info()
	if err != nil {
		return false, err
	}
	return executablePerm&inf.Mode().Perm() != 0, nil
}

func (*Executable) String() string {
	return "executable"
}
