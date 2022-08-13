package main

import (
	"os"
)

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
