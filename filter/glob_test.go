package filter_test

import (
	"strings"
	"testing"

	"github.com/komem3/fing/filter"
)

const globText = "baaabab"

var globTestPattern = []struct {
	patttern string
	match    bool
}{
	// match pattern
	{"baaabab", true},
	{"b***bab", true},
	{"*****ba*****ab", true},
	{"*ab", true},
	{"**ab", true},
	{"*baaabab", true},
	{"ba*", true},
	{"ba**", true},
	{"*ab*", true},
	{"**aaaba**", true},
	{"baaabab*", true},
	{"baa??ab", true},
	{"b*a?", true},
	{"b*b", true},
	{"?*", true},
	{"?a*ba?", true},
	{"??*??", true},
	// mismatch pattern
	{"a", false},
	{"b**a", false},
	{"**a", false},
	{"a*", false},
	{"*c*", false},
	{"baa", false},
	{"baaaba?b", false},
	{"bab", false},
	{"?", false},
	{"????", false},
}

func TestGlob_Match(t *testing.T) {
	t.Parallel()
	for _, tt := range globTestPattern {
		tt := tt
		t.Run(tt.patttern, func(t *testing.T) {
			t.Parallel()
			matcher := filter.NewGlob(tt.patttern)
			if match := filter.GlobMatch(matcher, globText); match != tt.match {
				t.Errorf("Match want %t, but got %t ", tt.match, match)
			}
		})
	}
}

func BenchmarkGlob_Match(b *testing.B) {
	repeata := strings.Repeat("a", 1000000)
	matcher := filter.NewGlob("*a*a*")
	for i := 0; i < b.N; i++ {
		if match := filter.GlobMatch(matcher, repeata); !match {
			b.Fatal("no match")
		}
	}
}

func BenchmarkGlob_Match_Equal(b *testing.B) {
	repeata := strings.Repeat("a", 1000000)
	matcher := filter.NewGlob(repeata)
	for i := 0; i < b.N; i++ {
		if match := filter.GlobMatch(matcher, repeata); !match {
			b.Fatal("no match")
		}
	}
}

func BenchmarkGlob_Match_Backward(b *testing.B) {
	repeata := strings.Repeat("a", 1000000)
	matcher := filter.NewGlob("a*")
	for i := 0; i < b.N; i++ {
		if match := filter.GlobMatch(matcher, repeata); !match {
			b.Fatal("no match")
		}
	}
}

func BenchmarkGlob_Match_Forward(b *testing.B) {
	repeata := strings.Repeat("a", 1000000)
	matcher := filter.NewGlob("*a")
	for i := 0; i < b.N; i++ {
		if match := filter.GlobMatch(matcher, repeata); !match {
			b.Fatal("no match")
		}
	}
}

func BenchmarkGlob_Match_ForwardBackwardk(b *testing.B) {
	repeata := strings.Repeat("a", 1000000)
	matcher := filter.NewGlob("*a*")
	for i := 0; i < b.N; i++ {
		if match := filter.GlobMatch(matcher, repeata); !match {
			b.Fatal("no match")
		}
	}
}
