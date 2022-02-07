package walk

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/komem3/fing/filter"
)

var concurrencyMax = runtime.NumCPU() * 2

type printType int

const (
	println printType = iota
	print0
)

const (
	defaultDirecotryBuffer = 1 << 7
)

type Walker struct {
	// matcher
	matcher filter.OrExp
	prunes  filter.OrExp

	// options
	IsDry     bool
	gitignore bool
	depth     int

	// result
	out     *bufio.Writer
	outerr  *bufio.Writer
	IsErr   bool
	targets directoryInfos
	writing sync.WaitGroup

	// print
	printType printType

	// concurrency control
	wg          sync.WaitGroup
	outmux      sync.Mutex
	dirMux      sync.Mutex
	concurrency chan struct{}

	fmt.Stringer
}

type direcotryInfo struct {
	path   string
	ignore *filter.Gitignore
}

type directoryInfos []*direcotryInfo

func (w *Walker) Walk(roots directoryInfos) {
	for _, r := range roots {
		f, err := os.Open(r.path)
		if err != nil {
			w.writeError(err)
			continue
		}
		entry, err := newEntry(f)
		f.Close()
		if err != nil {
			w.writeError(err)
			continue
		}
		w.walkFile(r.path, entry, nil)
	}

	depth := 1
	for len(roots) > 0 && (w.depth == -1 || depth <= w.depth) {
		w.targets = w.targets[:0]
		for i := range roots {
			w.wg.Add(1)
			w.concurrency <- struct{}{}
			go func(root *direcotryInfo) {
				w.walk(root.path, root.ignore)
				<-w.concurrency
				w.wg.Done()
			}(roots[i])
		}
		w.wg.Wait()

		if cap(roots) >= len(w.targets) {
			roots = roots[:len(w.targets)]
		} else {
			roots = make(directoryInfos, len(w.targets))
		}
		copy(roots, w.targets)

		depth++
	}
}

func (w *Walker) String() string {
	var s strings.Builder
	if w.gitignore {
		s.WriteString("ignore=true ")
	}
	if w.depth != -1 {
		fmt.Fprintf(&s, "maxdepth=%d ", w.depth)
	}
	if len(w.prunes) > 0 {
		fmt.Fprintf(&s, "prunes=[%s] ", w.prunes)
	}
	if len(w.matcher) > 0 {
		fmt.Fprintf(&s, "condition=[%s]", w.matcher)
	}
	return s.String()
}

func (w *Walker) walk(root string, gitignores *filter.Gitignore) {
	files, err := w.readDir(root)
	if err != nil {
		w.writeError(err)
		return
	}

	var newIgnore *filter.Gitignore
	if w.gitignore {
		ignoreFile := w.getIgnore(files)
		if ignoreFile != "" {
			newIgnore, err = filter.NewGitIgnore(root, ignoreFile)
			if err != nil {
				w.writeError(err)
				return
			}
		}
	}
	newIgnore = gitignores.Add(newIgnore)

	for i := range files {
		w.walkFile(filepath.Join(root, files[i].Name()), files[i], newIgnore)
	}
}

func (w *Walker) walkFile(path string, info fs.DirEntry, ignores *filter.Gitignore) {
	if ignores != nil {
		if match, _ := ignores.Match(path, info); match {
			return
		}
	}
	if info.IsDir() {
		match, err := w.prunes.Match(path, info)
		if err != nil {
			w.writeError(err)
			return
		}
		if len(w.prunes) > 0 && match {
			return
		}
		w.dirMux.Lock()
		if info.Name() == ".git" {
			w.targets = append(w.targets, &direcotryInfo{path: path})
		} else {
			w.targets = append(w.targets, &direcotryInfo{path: path, ignore: ignores})
		}
		w.dirMux.Unlock()
	}
	match, err := w.matcher.Match(path, info)
	if err != nil {
		w.writeError(err)
		return
	}
	if match {
		w.writeFile(path, info)
	}
}

func (w *Walker) writeError(err error) {
	w.IsErr = true
	w.writing.Add(1)
	go func() {
		w.outmux.Lock()
		if _, err := w.outerr.WriteString(err.Error() + "\n"); err != nil {
			log.Printf("[ERROR] %v", err)
		}
		w.outmux.Unlock()
		w.writing.Done()
	}()
}

func (w *Walker) writeFile(path string, _ fs.DirEntry) {
	w.writing.Add(1)
	go func() {
		w.outmux.Lock()
		switch w.printType {
		case println:
			if _, err := w.out.WriteString(path + "\n"); err != nil {
				log.Printf("[ERROR] %v", err)
			}
		case print0:
			if _, err := w.out.WriteString(path + "\x00"); err != nil {
				log.Printf("[ERROR] %v", err)
			}
		}
		w.outmux.Unlock()
		w.writing.Done()
	}()
}

func (w *Walker) Wait() {
	w.writing.Wait()
}

func (w *Walker) Flush() {
	w.outmux.Lock()
	if err := w.out.Flush(); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	if err := w.outerr.Flush(); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	w.outmux.Unlock()
}

func (w *Walker) readDir(dir string) (ds []fs.DirEntry, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	ds, err = f.ReadDir(-1)
	f.Close()
	return
}

func (*Walker) getIgnore(files []fs.DirEntry) string {
	for i := range files {
		if files[i].Name() == ".gitignore" {
			return files[i].Name()
		}
	}
	return ""
}

func (d directoryInfos) String() string {
	paths := make([]string, 0, len(d))
	for _, p := range d {
		paths = append(paths, p.path)
	}
	return strings.Join(paths, ", ")
}
