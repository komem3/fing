package main

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/komem3/fing/walk"
)

func run(args []string, stdout, stderr io.Writer) error {
	out := bufio.NewWriterSize(stdout, 0)
	outerr := bufio.NewWriterSize(stderr, 0)
	walker, paths, err := walk.NewWalkerFromArgs(args, out, outerr)
	if err != nil {
		return err
	}
	if walker.IsDry {
		fmt.Fprintf(stdout, "targets=[%s] %s\n", paths, walker)
		return nil
	}

	ch := make(chan struct{}, 1)
	ticker := time.NewTicker(time.Millisecond)
	go func() {
		walker.Walk(paths)
		walker.Wait()
		ch <- struct{}{}
	}()
	for end := false; !end; {
		select {
		case <-ch:
			end = true
		case <-ticker.C:
			walker.Flush()
		}
	}
	walker.Flush()

	if walker.IsErr {
		return fmt.Errorf("error occurred")
	}
	return nil
}
