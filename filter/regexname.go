package filter

import (
	"fmt"
	"io/fs"
	"regexp"
)

type RegexName struct {
	reg *regexp.Regexp
}

var _ FileExp = (*RegexName)(nil)

func NewRegexName(pattern string) (*RegexName, error) {
	reg, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return nil, err
	}
	return &RegexName{
		reg: reg,
	}, nil
}

func NewIRegexName(pattern string) (*RegexName, error) {
	return NewRegexName("(?i)" + pattern)
}

func (r *RegexName) Match(_ string, info fs.DirEntry) bool {
	return r.reg.MatchString(info.Name())
}

func (r *RegexName) String() string {
	return fmt.Sprintf("regex_name(%s)", r.reg.String())
}
