package filter

import (
	"fmt"
	"io/fs"
	"regexp"
)

type Regex struct {
	reg *regexp.Regexp
}

var _ FileExp = (*Regex)(nil)

func NewRegex(pattern string) (*Regex, error) {
	reg, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return nil, err
	}
	return &Regex{
		reg: reg,
	}, nil
}

func NewIRegex(pattern string) (*Regex, error) {
	return NewRegex("(?i)" + pattern)
}

func (r *Regex) Match(path string, _ fs.DirEntry) bool {
	return r.reg.MatchString(path)
}

func (r *Regex) String() string {
	return fmt.Sprintf("regex(%s)", r.reg.String())
}
