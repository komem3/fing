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
	matcher       filter.OrExp
	prunes        filter.OrExp
	excludeIgnore filter.OrExp

	// options
	IsDry     bool
	gitignore bool
	depth     int

	// result
	out     *bufio.Writer
	outerr  *bufio.Writer
	IsErr   bool
	targets entryInfos
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

type entryInfo struct {
	path   string
	ignore *filter.Gitignore
	info   fs.DirEntry
}

type entryInfos []*entryInfo

func (w *Walker) Walk(roots []string) {
	entries := make([]*entryInfo, 0, len(roots))
	for _, r := range roots {
		f, err := os.Open(r)
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
		entries = append(entries, &entryInfo{path: r, info: entry})
	}

	for depth := 0; len(entries) > 0 && (w.depth == -1 || depth <= w.depth); depth++ {
		w.targets = w.targets[:0]
		for i := range entries {
			w.wg.Add(1)
			w.concurrency <- struct{}{}
			go func(entry *entryInfo) {
				w.walk(entry)
				<-w.concurrency
				w.wg.Done()
			}(entries[i])
		}
		w.wg.Wait()

		if cap(entries) >= len(w.targets) {
			entries = entries[:len(w.targets)]
		} else {
			entries = make(entryInfos, len(w.targets))
		}
		copy(entries, w.targets)
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

func (w *Walker) walk(entry *entryInfo) {
	if entry.ignore != nil {
		var (
			match bool
			err   error
		)
		if len(w.excludeIgnore) > 0 {
			match, err = w.excludeIgnore.Match(entry.path, entry.info)
			if err != nil {
				w.writeError(err)
				return
			}
		}
		if !match {
			if match, _ := entry.ignore.Match(entry.path, entry.info); match {
				return
			}
		}
	}
	match, err := w.matcher.Match(entry.path, entry.info)
	if err != nil {
		w.writeError(err)
		return
	}
	if match {
		w.writeFile(entry.path, entry.info)
	}

	if entry.info.IsDir() {
		w.scanDir(entry)
	}
}

func (w *Walker) scanDir(entry *entryInfo) {
	match, err := w.prunes.Match(entry.path, entry.info)
	if err != nil {
		w.writeError(err)
		return
	}
	if len(w.prunes) > 0 && match {
		return
	}

	files, err := w.readDir(entry.path)
	if err != nil {
		w.writeError(err)
		return
	}

	var newIgnore *filter.Gitignore
	if w.gitignore {
		ignoreFile := w.getIgnore(files)
		if ignoreFile != "" {
			newIgnore, err = filter.NewGitIgnore(entry.path, ignoreFile)
			if err != nil {
				w.writeError(err)
				return
			}
		}
	}
	newIgnore = entry.ignore.Add(newIgnore)

	for _, f := range files {
		w.dirMux.Lock()
		if entry.info.Name() == ".git" {
			w.targets = append(w.targets, &entryInfo{path: filepath.Join(entry.path, f.Name()), info: f})
		} else {
			w.targets = append(w.targets, &entryInfo{path: filepath.Join(entry.path, f.Name()), info: f, ignore: newIgnore})
		}
		w.dirMux.Unlock()
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

func (d entryInfos) String() string {
	paths := make([]string, 0, len(d))
	for _, p := range d {
		paths = append(paths, p.path)
	}
	return strings.Join(paths, ", ")
}
