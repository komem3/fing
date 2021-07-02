package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "\n", err)
		os.Exit(1)
	}
}
