package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/komem3/fing/walk"
)

func run(args []string, stdout, stderr io.Writer) (status int) {
	walker, paths, err := walk.NewWalkerFromArgs(args, stdout, stderr)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return 1
	}
	if walker.IsDry {
		fmt.Fprintf(stdout, "targets=[%s] %s\n", strings.Join(paths, ", "), walker)
		return 0
	}

	walker.Walk(paths)
	if walker.IsErr {
		return 1
	}
	return 0
}
