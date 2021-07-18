package filter

import (
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"
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

	runeNums int

	matchPattern matchPattern
	matchIndex   int

	fmt.Stringer
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
		runeNums := utf8.RuneCountInString(pattern)
		return &glob{
			pattern:      pattern,
			matchPattern: globPattern,
			ddlPool: sync.Pool{
				New: func() interface{} {
					ddl := make([][]bool, 0, runeNums)
					return &ddl
				},
			},
			runeNums: runeNums,
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

	str := []rune(s)
	pool := g.ddlPool.Get().(*[][]bool)
	dp := *pool
	dp = dp[:cap(dp)]

	var (
		firstMatch int
		dpIndex    int
	)
	for pIndex, r := range g.pattern {
		if cap(dp[dpIndex]) > 0 {
			dp[dpIndex] = dp[dpIndex][:0]
		} else if dpIndex > 2 {
			dp[dpIndex] = dp[dpIndex-2][:0]
		} else {
			dp[dpIndex] = make([]bool, 0, defaultBuf)
		}
		var match bool
		switch g.pattern[dpIndex] {
		case '*':
			if dpIndex == len(g.pattern)-1 {
				*pool = dp[:0]
				g.ddlPool.Put(pool)
				return true
			}
			for j := 0; j < len(s)-firstMatch; j++ {
				dp[dpIndex] = append(dp[dpIndex], true)
			}
		case '?':
			for arrayj, strj := 0, firstMatch; strj < len(str); arrayj, strj = arrayj+1, strj+1 {
				m := (dpIndex == 0 && arrayj == 0) || (dpIndex != 0 && arrayj != 0 && dp[dpIndex-1][arrayj-1])
				if !match {
					if !m {
						continue
					}
					match = true
					firstMatch = strj
				}
				dp[dpIndex] = append(dp[dpIndex], m)
			}
		default:
			if dpIndex == 0 {
				if r != str[0] {
					*pool = dp[:0]
					g.ddlPool.Put(pool)
					return false
				}
			}
			if dpIndex != 0 && g.pattern[pIndex-1] == '*' {
				for arrayj, strj := 0, firstMatch; strj < len(str); arrayj, strj = arrayj+1, strj+1 {
					m := (dpIndex == 0 || dp[dpIndex-1][arrayj]) && r == str[strj]
					if !match {
						if !m {
							continue
						}
						match = true
						firstMatch = strj
					}
					dp[dpIndex] = append(dp[dpIndex], m)
				}
			} else {
				for arrayj, strj := 0, firstMatch; strj < len(str); arrayj, strj = arrayj+1, strj+1 {
					var m bool
					if dpIndex == 0 {
						m = arrayj == 0 && r == str[strj]
					} else {
						m = arrayj != 0 && dp[dpIndex-1][arrayj-1] && r == str[strj]
					}
					if !match {
						if !m {
							continue
						}
						match = true
						firstMatch = strj
					}
					dp[dpIndex] = append(dp[dpIndex], m)
				}
			}
			if len(dp[dpIndex]) == 0 {
				*pool = dp[:0]
				g.ddlPool.Put(pool)
				return false
			}
		}
		dpIndex++
	}
	*pool = dp[:0]
	g.ddlPool.Put(pool)
	return len(dp[g.runeNums-1])+firstMatch == len(str) && dp[g.runeNums-1][len(dp[g.runeNums-1])-1]
}

func (g *glob) String() string {
	return g.pattern
}
