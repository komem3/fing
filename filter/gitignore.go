package filter

import (
	"io/fs"
)

type Gitignore struct {
	PathMatchers []*Path
}

var _ FileExp = (*Gitignore)(nil)

func (g *Gitignore) Match(path string, info fs.DirEntry) (bool, error) {
	var match bool
	for i := range g.PathMatchers {
		m, err := g.PathMatchers[i].Match(path, info)
		if err != nil {
			return false, err
		}
		if m {
			match = g.PathMatchers[i].pathType == normalPathType
		}
	}
	return match, nil
}

func (g *Gitignore) Add(src *Gitignore) *Gitignore {
	if src == nil {
		return g
	}
	if g == nil {
		return src
	}
	dst := make([]*Path, len(src.PathMatchers)+len(g.PathMatchers))
	copy(dst, src.PathMatchers)
	copy(dst[len(src.PathMatchers):], g.PathMatchers)
	return &Gitignore{PathMatchers: dst}
}
