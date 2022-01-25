package filter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

var separator = string(filepath.Separator)

type Gitignore struct {
	PathMatchers []gitignore.Pattern
}

var _ FileExp = (*Gitignore)(nil)

func NewGitIgnore(dir string, file string) (*Gitignore, error) {
	domain := strings.Split(dir, separator)
	if len(domain) > 0 && domain[0] == "." {
		domain = domain[1:]
	}
	buf, err := os.ReadFile(filepath.Join(dir, file))
	if err != nil {
		return nil, err
	}
	var ignores []gitignore.Pattern
	reader := bufio.NewReader(bytes.NewReader(buf))
	for {
		b, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if len(b) == 0 ||
			b[0] == '#' {
			continue
		}
		ignores = append(ignores, gitignore.ParsePattern(string(b), domain))
	}
	return &Gitignore{PathMatchers: ignores}, nil
}

func (g *Gitignore) Match(path string, info fs.DirEntry) (bool, error) {
	var match bool
	splitPath := strings.Split(path, separator)
	for i := range g.PathMatchers {
		if m := g.PathMatchers[i].Match(splitPath, info.IsDir()); m > gitignore.NoMatch {
			match = m == gitignore.Exclude
		}
	}
	return match, nil
}

func (g *Gitignore) Add(src *Gitignore) *Gitignore {
	if (g == nil || len(g.PathMatchers) == 0) &&
		(src == nil || len(src.PathMatchers) == 0) {
		return nil
	}
	if src == nil || len(src.PathMatchers) == 0 {
		dst := make([]gitignore.Pattern, len(g.PathMatchers))
		copy(dst, g.PathMatchers)
		return &Gitignore{PathMatchers: dst}
	}
	if g == nil || len(g.PathMatchers) == 0 {
		dst := make([]gitignore.Pattern, len(src.PathMatchers))
		copy(dst, src.PathMatchers)
		return &Gitignore{PathMatchers: dst}
	}
	dst := make([]gitignore.Pattern, len(src.PathMatchers)+len(g.PathMatchers))
	copy(dst, g.PathMatchers)
	copy(dst[len(g.PathMatchers):], src.PathMatchers)
	return &Gitignore{PathMatchers: dst}
}
