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

func (e OrExp) Match(path string, info fs.DirEntry) bool {
	if len(e) == 0 {
		return true
	}
	for _, filters := range e {
		if filters.Match(path, info) {
			return true
		}
	}
	return false
}

func (e AndExp) Match(path string, info fs.DirEntry) bool {
	for _, filters := range e {
		if !filters.Match(path, info) {
			return false
		}
	}
	return true
}

func (n *NotExp) Match(path string, info fs.DirEntry) bool {
	return !n.filter.Match(path, info)
}

func (e OrExp) String() string {
	var buf strings.Builder
	for i, f := range e {
		if i == 0 {
			fmt.Fprintf(&buf, "%s", f)
			continue
		}
		fmt.Fprintf(&buf, " + %s", f)
	}
	return buf.String()
}

func (e AndExp) String() string {
	var buf strings.Builder
	for i, f := range e {
		if i == 0 {
			fmt.Fprintf(&buf, "%s", f)
			continue
		}
		fmt.Fprintf(&buf, " * %s", f)
	}
	return buf.String()
}

func (e *NotExp) String() string {
	return fmt.Sprintf("not %s", e.filter)
}
