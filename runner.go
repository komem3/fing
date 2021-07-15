package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/komem3/fing/walk"
)

// If there are many allocates, the execution speed will decrease.
// Therefore, the initial value is increased.
const initBufferSize = 1 << 10

func run(args []string, stdout, stderr io.Writer) error {
	out := bufio.NewWriterSize(stdout, initBufferSize)
	outerr := bufio.NewWriterSize(stderr, initBufferSize)
	walker, paths, err := walk.NewWalkerFromArgs(args, out, outerr)
	if err != nil {
		return err
	}
	if walker.IsDry {
		fmt.Fprintf(stdout, "%s\n", walker)
		return nil
	}
	walker.Walk(paths)

	if err := out.Flush(); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	if err := outerr.Flush(); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	if walker.IsErr {
		return fmt.Errorf("error occurred")
	}
	return nil
}
