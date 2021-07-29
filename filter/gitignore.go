package filter

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type Gitignore struct {
	PathMatchers []gitignore.Pattern
}

var _ FileExp = (*Gitignore)(nil)

func (g *Gitignore) Match(path string, info fs.DirEntry) (bool, error) {
	var match bool
	for i := range g.PathMatchers {
		if m := g.PathMatchers[i].Match(strings.Split(path, string(filepath.Separator)), info.IsDir()); m != gitignore.NoMatch {
			match = m == gitignore.Exclude
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
	dst := make([]gitignore.Pattern, len(src.PathMatchers)+len(g.PathMatchers))
	copy(dst, g.PathMatchers)
	copy(dst[len(g.PathMatchers):], src.PathMatchers)
	return &Gitignore{PathMatchers: dst}
}
