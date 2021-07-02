package filter

import (
	"strings"
	"sync"
)

type matchPattern int

const (
	equalPatten matchPattern = iota
	globPattern
	backwardPattern
	forwardPattern
	forwardBackwardPattern
)

type glob struct {
	pattern string
	ddlPool sync.Pool

	matchPattern matchPattern
	matchIndex   int
}

const defaultBuf = 1 << 10

func newGlob(pattern string) *glob {
	var (
		matchp     matchPattern
		matchIndex int
		lastIndex  int
	)
loopend:
	for i, r := range pattern {
		switch {
		case r == '?':
			matchp = globPattern
			break loopend

		// backward
		case i == 0 && r == '*':
			matchp = backwardPattern
			matchIndex = len(pattern) - 1

		case matchp == backwardPattern && r == '*':
			if pattern[i-1] == '*' {
				matchIndex--
				continue
			}
			matchIndex = len(pattern) - matchIndex
			lastIndex = i
			matchp = forwardBackwardPattern

		// forwardBackward
		case matchp == forwardBackwardPattern && r != '*':
			matchp = globPattern
			break loopend

		// forward
		case matchp == equalPatten && r == '*':
			matchp = forwardPattern
			matchIndex = i

		case matchp == forwardPattern && r != '*':
			matchp = globPattern
			break loopend
		}
	}
	if matchp == globPattern {
		return &glob{
			pattern:      pattern,
			matchPattern: globPattern,
			ddlPool: sync.Pool{
				New: func() interface{} {
					ddl := make([][]bool, 0, len(pattern))
					return &ddl
				},
			},
		}
	}
	if matchp == forwardBackwardPattern {
		return &glob{
			pattern:      pattern[matchIndex:lastIndex],
			matchPattern: matchp,
		}
	}
	return &glob{
		pattern:      pattern,
		matchPattern: matchp,
		matchIndex:   matchIndex,
	}
}

func (g *glob) match(s string) bool {
	switch g.matchPattern {
	case equalPatten:
		return g.pattern == s
	case forwardPattern:
		return len(s) >= g.matchIndex && s[:g.matchIndex] == g.pattern[:g.matchIndex]
	case backwardPattern:
		return len(s) >= g.matchIndex && s[len(s)-g.matchIndex:] == g.pattern[len(g.pattern)-g.matchIndex:]
	case forwardBackwardPattern:
		return strings.Contains(s, g.pattern)
	}

	pool := g.ddlPool.Get().(*[][]bool)
	ddl := *pool
	ddl = ddl[:cap(ddl)]

	var firstMatch int
	for i := range g.pattern {
		if cap(ddl[i]) > 0 {
			ddl[i] = ddl[i][:0]
		} else if i > 2 {
			ddl[i] = ddl[i-2][:0]
		} else {
			ddl[i] = make([]bool, 0, defaultBuf)
		}
		var match bool
		switch g.pattern[i] {
		case '*':
			if i == len(g.pattern)-1 {
				*pool = ddl[:0]
				g.ddlPool.Put(pool)
				return true
			}
			for j := 0; j < len(s)-firstMatch; j++ {
				ddl[i] = append(ddl[i], true)
			}
		case '?':
			for arrayj, strj := 0, firstMatch; strj < len(s); arrayj, strj = arrayj+1, strj+1 {
				m := (i == 0 && arrayj == 0) || (i != 0 && arrayj != 0 && ddl[i-1][arrayj-1])
				if !match {
					if !m {
						continue
					}
					match = true
					firstMatch = strj
				}
				ddl[i] = append(ddl[i], m)
			}
		default:
			if i == 0 {
				if g.pattern[i] != s[0] {
					*pool = ddl[:0]
					g.ddlPool.Put(pool)
					return false
				}
			}
			if i != 0 && g.pattern[i-1] == '*' {
				for arrayj, strj := 0, firstMatch; strj < len(s); arrayj, strj = arrayj+1, strj+1 {
					m := (i == 0 || ddl[i-1][arrayj]) && g.pattern[i] == s[strj]
					if !match {
						if !m {
							continue
						}
						match = true
						firstMatch = strj
					}
					ddl[i] = append(ddl[i], m)
				}
			} else {
				for arrayj, strj := 0, firstMatch; strj < len(s); arrayj, strj = arrayj+1, strj+1 {
					var m bool
					if i == 0 {
						m = arrayj == 0 && g.pattern[i] == s[strj]
					} else {
						m = arrayj != 0 && ddl[i-1][arrayj-1] && g.pattern[i] == s[strj]
					}
					if !match {
						if !m {
							continue
						}
						match = true
						firstMatch = strj
					}
					ddl[i] = append(ddl[i], m)
				}
			}
			if len(ddl[i]) == 0 {
				*pool = ddl[:0]
				g.ddlPool.Put(pool)
				return false
			}
		}
	}
	*pool = ddl[:0]
	g.ddlPool.Put(pool)
	return len(ddl[len(g.pattern)-1])+firstMatch == len(s) && ddl[len(g.pattern)-1][len(ddl[len(g.pattern)-1])-1]
}
