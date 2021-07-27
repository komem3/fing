package filter

import (
	"fmt"
	"io/fs"
	"strconv"
)

type CmpOption int

const (
	EqualCmpOption CmpOption = iota
	GreaterCmpOption
	LessCmpOption
)

type Size struct {
	Size int64
	Opt  CmpOption
}

var _ FileExp = (*Size)(nil)

func NewSize(str string) (*Size, error) {
	if len(str) == 0 {
		return nil, fmt.Errorf("missing argument of size")
	}
	var (
		opt CmpOption
		s   = str[:]
	)
	switch s[0] {
	case '+':
		opt = GreaterCmpOption
		if len(s) == 1 {
			return nil, fmt.Errorf("%s is invalid size argument", str)
		}
		s = s[1:]
	case '-':
		opt = LessCmpOption
		s = s[1:]
	default:
		opt = EqualCmpOption
	}

	var size int64
	switch s[len(s)-1] {
	case 'c':
		psize, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		size = psize
	case 'k':
		psize, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		size = psize * 1024
	case 'M':
		psize, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		size = psize * 1024 * 1024
	case 'G':
		psize, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		size = psize * 1024 * 1024 * 1024
	default:
		return nil, fmt.Errorf("%c is invalid unit of size", s[len(s)-1])
	}
	return &Size{
		Size: size,
		Opt:  opt,
	}, nil
}

func (s *Size) Match(_ string, entry fs.DirEntry) (bool, error) {
	info, err := entry.Info()
	if err != nil {
		return false, err
	}
	switch s.Opt {
	case EqualCmpOption:
		return info.Size() == int64(s.Size), nil
	case GreaterCmpOption:
		return info.Size() > int64(s.Size), nil
	case LessCmpOption:
		return info.Size() < int64(s.Size), nil
	}
	panic("invalid compare option")
}
