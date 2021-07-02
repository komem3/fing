package main

import (
	"io"
	"testing"
)

func BenchmarkRun_IName(b *testing.B) {
	args := []string{"fing", "testdata", "-iname", "*.png"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := run(args, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRun_IPath(b *testing.B) {
	args := []string{"fing", "testdata", "-ipath", "*.png"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := run(args, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRun_Regex(b *testing.B) {
	args := []string{"fing", "testdata", "-iregex", `.*\.png`}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := run(args, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
