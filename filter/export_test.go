package filter

func (g *Gitignore) Len() int {
	if g == nil {
		return 0
	}
	return len(g.PathMatchers)
}
