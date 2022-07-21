package filter

import (
	"fmt"
	"io/fs"
	"strings"
)

type OrExp []FileExp

type AndExp []FileExp

type NotExp struct {
	filter FileExp
}

type AlwasyExp bool

var (
	_ FileExp = (OrExp)(nil)
	_ FileExp = (AndExp)(nil)
	_ FileExp = (*NotExp)(nil)
)

func NewNotExp(f FileExp) *NotExp {
	return &NotExp{f}
}

func NewOrExp(f FileExp) OrExp {
	or := make(OrExp, 1)
	or[0] = f
	return or
}

func (e OrExp) Match(path string, info fs.DirEntry) (bool, error) {
	if len(e) == 0 {
		return true, nil
	}
	for _, filters := range e {
		match, err := filters.Match(path, info)
		if err != nil {
			return false, err
		}
		if match {
			return match, nil
		}
	}
	return false, nil
}

func (e AndExp) Match(path string, info fs.DirEntry) (bool, error) {
	for _, filters := range e {
		match, err := filters.Match(path, info)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (n *NotExp) Match(path string, info fs.DirEntry) (bool, error) {
	match, err := n.filter.Match(path, info)
	if err != nil {
		return false, err
	}
	return !match, nil
}

func (e OrExp) String() string {
	var buf strings.Builder
	for i, f := range e {
		if i == 0 {
			fmt.Fprintf(&buf, "%s", f)
			continue
		}
		fmt.Fprintf(&buf, " || %s", f)
	}
	return buf.String()
}

func (a AlwasyExp) Match(path string, info fs.DirEntry) (bool, error) {
	return bool(a), nil
}

func (e AndExp) String() string {
	var buf strings.Builder
	for i, f := range e {
		if i == 0 {
			fmt.Fprintf(&buf, "%s", f)
			continue
		}
		fmt.Fprintf(&buf, " && %s", f)
	}
	return buf.String()
}

func (e *NotExp) String() string {
	return fmt.Sprintf("not %s", e.filter)
}

func (a AlwasyExp) String() string {
	return fmt.Sprintf("%t", bool(a))
}
