package walk

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/komem3/fing/filter"
)

const openFileMax = 1 << 9

const defaultIgnoreBuffer = 1 << 7

type Walker struct {
	// matcher
	matcher filter.OrExp
	prunes  filter.OrExp

	// options
	gitignore bool

	// result
	out    io.Writer
	outerr io.Writer
	IsErr  bool

	// concurrency control
	wg        sync.WaitGroup
	mux       sync.Mutex
	openFiles chan struct{}
}

func (w *Walker) Wait() {
	w.wg.Wait()
}

func (w *Walker) Walk(root string) {
	w.wg.Add(1)
	go func() {
		w.walk(root, nil)
		w.wg.Done()
	}()
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
			newIgnore, err = w.extractGitignore(root, ignoreFile)
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
		if ignores.Match(path, info) {
			return
		}
	}
	if info.IsDir() {
		if len(w.prunes) > 0 && w.prunes.Match(path, info) {
			return
		}
		w.wg.Add(1)
		go func() {
			if info.Name() == ".git" {
				w.walk(path, nil)
			} else {
				w.walk(path, ignores)
			}
			w.wg.Done()
		}()
	}
	if w.matcher.Match(path, info) {
		w.writeFile(path, info)
	}
}

func (w *Walker) writeError(err error) {
	w.IsErr = true
	w.mux.Lock()
	if _, err := fmt.Fprintln(w.outerr, err.Error()); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	w.mux.Unlock()
}

func (w *Walker) writeFile(path string, _ fs.DirEntry) {
	w.mux.Lock()
	if _, err := fmt.Fprintln(w.out, path); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	w.mux.Unlock()
}

func (w *Walker) readDir(dir string) (ds []fs.DirEntry, err error) {
	w.openFiles <- struct{}{}
	f, err := os.Open(dir)
	if err != nil {
		<-w.openFiles
		return nil, err
	}
	ds, err = f.ReadDir(-1)
	f.Close()
	<-w.openFiles
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

func (w *Walker) extractGitignore(root, path string) (ignore *filter.Gitignore, err error) {
	w.openFiles <- struct{}{}
	buf, err := os.ReadFile(filepath.Join(root, path))
	<-w.openFiles
	if err != nil {
		return nil, err
	}
	ignore = &filter.Gitignore{
		PathMatchers: make([]*filter.Path, 0, defaultIgnoreBuffer),
	}
	for _, file := range strings.Split(string(buf), "\n") {
		for len(file) > 2 && file[len(file)-2] == '*' && file[len(file)-1] == '*' {
			file = file[:len(file)-1]
		}
		if len(file) > 2 && file[len(file)-2:] == "/*" {
			file = file[:len(file)-2]
		}
		file = strings.TrimRight(file, "/")
		if len(file) == 0 ||
			strings.HasPrefix(file, "#") ||
			(file[0] == '!' && len(file) == 1) ||
			(file[0] == '/' && len(file) == 1) ||
			(file[0] == '!' && file[1] == '/' && len(file) == 2) {
			continue
		}
		if file[0] == '!' {
			if file[1] == '*' {
				ignore.PathMatchers = append(ignore.PathMatchers, filter.NewNotPath(file[1:]))
				continue
			}
			if file[1] != '/' {
				ignore.PathMatchers = append(ignore.PathMatchers, filter.NewNotPath(filepath.Join("*", file[1:])))
			}
			ignore.PathMatchers = append(ignore.PathMatchers, filter.NewNotPath(filepath.Join(root, file[1:])))
			continue
		}
		if file[0] == '*' {
			ignore.PathMatchers = append(ignore.PathMatchers, filter.NewPath(file))
			continue
		}
		if file[0] != '/' {
			ignore.PathMatchers = append(ignore.PathMatchers, filter.NewPath(filepath.Join("*", file)))
		}
		ignore.PathMatchers = append(ignore.PathMatchers, filter.NewPath(filepath.Join(root, file)))
	}
	return ignore, nil
}
