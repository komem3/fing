package filter

var NewGlob = newGlob

var GlobMatch = (*glob).match

func (g *Gitignore) Len() int {
	if g == nil {
		return 0
	}
	return len(g.PathMatchers)
}
