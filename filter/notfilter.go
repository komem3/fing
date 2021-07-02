package filter

import "io/fs"

type NotFilter struct {
	filter FileExp
}

var _ FileExp = (*NotFilter)(nil)

func NewNotFilter(f FileExp) *NotFilter {
	return &NotFilter{f}
}

func (n *NotFilter) Match(path string, info fs.DirEntry) (bool, error) {
	m, err := n.filter.Match(path, info)
	return !m, err
}
