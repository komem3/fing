package walk

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/komem3/fing/filter"
)

var concurrencyNum = runtime.NumCPU() * 8

type printType int

const (
	println printType = iota
	print0
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
	ignoreErr bool

	// result
	out         *bufio.Writer
	outerr      io.Writer
	flushTick   *time.Ticker
	IsErr       bool
	directories entryInfos

	// print
	printType printType

	// concurrency control
	writingMutex sync.Mutex
	dirMutex     sync.Mutex

	fmt.Stringer
}

type entryInfo struct {
	path   string
	ignore *filter.Gitignore
	info   fs.DirEntry
}

type entryInfos []*entryInfo

func (w *Walker) Walk(roots []string) {
	w.flushTick = time.NewTicker(time.Millisecond)
	defer w.flushTick.Stop()

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

		w.checkEntry(&entryInfo{path: r, info: entry})
	}

	var wg sync.WaitGroup
	dirsChans := make([]chan entryInfos, concurrencyNum)
	for i := 0; i < concurrencyNum; i++ {
		i := i
		dirsChans[i] = make(chan entryInfos)
		go func() {
			for dirs := range dirsChans[i] {
				for _, dir := range dirs {
					w.scanDir(dir)
				}
				wg.Done()
			}
		}()
		defer func() {
			close(dirsChans[i])
		}()
	}

	var entries entryInfos
	for depth := 1; len(w.directories) > 0 && (w.depth == -1 || depth <= w.depth); depth++ {
		if cap(entries) >= len(w.directories) {
			entries = entries[:len(w.directories)]
		} else {
			entries = make(entryInfos, len(w.directories))
		}
		copy(entries, w.directories)
		w.directories = w.directories[:0]

		if len(entries) < concurrencyNum {
			wg.Add(1)
			dirsChans[0] <- entries
		} else {
			chunkSize := len(entries) / concurrencyNum
			for i := 0; i < concurrencyNum; i++ {
				wg.Add(1)
				if i == concurrencyNum-1 {
					dirsChans[i] <- entries[i*chunkSize:]
				} else {
					dirsChans[i] <- entries[i*chunkSize : (i+1)*chunkSize]
				}
			}
		}
		wg.Wait()
	}

	if err := w.out.Flush(); err != nil {
		log.Printf("[ERROR] %v", err)
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

func (w *Walker) checkEntry(entry *entryInfo) {
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
		w.dirMutex.Lock()
		w.directories = append(w.directories, entry)
		w.dirMutex.Unlock()
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
		if entry.info.Name() == ".git" {
			w.checkEntry(&entryInfo{path: filepath.Join(entry.path, f.Name()), info: f})
		} else {
			w.checkEntry(&entryInfo{path: filepath.Join(entry.path, f.Name()), info: f, ignore: newIgnore})
		}
	}
}

func (w *Walker) writeError(err error) {
	if w.ignoreErr {
		return
	}
	w.IsErr = true
	w.writingMutex.Lock()
	if _, err := w.outerr.Write([]byte(err.Error() + "\n")); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	w.writingMutex.Unlock()
}

func (w *Walker) writeFile(path string, _ fs.DirEntry) {
	w.writingMutex.Lock()
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

	select {
	case <-w.flushTick.C:
		if err := w.out.Flush(); err != nil {
			log.Printf("[ERROR] %v", err)
		}
	default:
	}
	w.writingMutex.Unlock()
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
