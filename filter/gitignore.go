package filter

import (
	"io/fs"
)

type Gitignore struct {
	PathMatchers []*Path
}

var _ FileExp = (*Gitignore)(nil)

func (g *Gitignore) Match(path string, info fs.DirEntry) bool {
	var match bool
	for i := range g.PathMatchers {
		if g.PathMatchers[i].Match(path, info) {
			match = g.PathMatchers[i].pathType == normalPathType
		}
	}
	return match
}

func (g *Gitignore) Add(src *Gitignore) *Gitignore {
	if src == nil {
		return g
	}
	if g == nil {
		return src
	}
	dst := make([]*Path, len(src.PathMatchers)+len(g.PathMatchers))
	copy(dst, g.PathMatchers)
	copy(dst[len(g.PathMatchers):], src.PathMatchers)
	return &Gitignore{PathMatchers: dst}
}
