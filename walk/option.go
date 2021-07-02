package walk

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/komem3/fing/filter"
)

const defaultMakeLen = 1 << 2

var (
	isGC         bool
	isNot        bool
	isNotProcess bool
)

type boolFunc func(bool)

func (boolFunc) String() string { return "false" }

func (boolFunc) IsBoolFlag() bool { return true }

func (i boolFunc) Set(s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	i(b)
	return nil
}

var Usage = `
Usage: fing [staring-point...] [options|expression]

Fing is a simple and fast file finder.

options are:
  -I
    Ignore files in .gitignore.
  -timeout
    Search time limit. (default 10s)
  -enable_gc
    Garbage collection is enable option.
    Default is disabled to speed up execution.

expression are:
  -iname string
		Like -name, but the match is case insensitive.
  -ipath string
		Like -path, but the match is case insensitive.
  -iregex string
		Like -regex, but the match is case insensitive.
  -name string
    Search for files using wildcard expressions.
    This option match only to file name.
  -not
   True if next expression false.
  -not( expression -not)
   True if next expression false in parentheses.
   No support nested parenthese.
  -path string
    Search for files using wildcard expressions.
    This option match to file path.
  -prune
    Prunes files and directory that match before expressions.
  -regex string
    Search for files using regular expressions.
    This option match to file path.
    Unlike find, this is a backward match.
  -type string
    File is type.
    Support file(f), directory(d), named piep(p) and socket(s).
`

func NewWalkerFromArgs(args []string, out, outerr io.Writer) (*Walker, []string, error) {
	walker := &Walker{
		filters:   make([]filter.FileExp, 0, defaultMakeLen),
		prunes:    make([]filter.FileExp, 0, defaultMakeLen),
		out:       out,
		outerr:    outerr,
		openFiles: make(chan struct{}, openFileMax),
	}

	flag := flag.NewFlagSet(args[0], flag.ExitOnError)
	flag.Usage = func() { fmt.Fprint(os.Stderr, Usage) }
	{
		// options
		flag.BoolVar(&walker.gitignore, "I", false, "")
		flag.DurationVar(&walker.timeout, "timeout", time.Second*10, "")
		flag.BoolVar(&isGC, "enable_gc", false, "")
	}
	{
		// expression
		flag.BoolVar(&isNot, "not", false, "")
		flag.Var(boolFunc(func(b bool) {
			isNotProcess = b
		}), "not(", "")
		flag.Var(boolFunc(func(b bool) {
			isNotProcess = !b
		}), "not)", "")
		flag.Func("name", "", func(s string) error {
			walker.filters = append(walker.filters, toFilter(filter.NewFileName(s)))
			return nil
		})
		flag.Func("iname", "", func(s string) error {
			walker.filters = append(walker.filters, toFilter(filter.NewIFileName(s)))
			return nil
		})
		flag.Func("path", "", func(s string) error {
			walker.filters = append(walker.filters, toFilter(filter.NewPath(s)))
			return nil
		})
		flag.Func("ipath", "", func(s string) error {
			walker.filters = append(walker.filters, toFilter(filter.NewIPath(s)))
			return nil
		})
		flag.Func("regex", "", func(s string) error {
			f, err := filter.NewRegex(s)
			if err != nil {
				return err
			}
			walker.filters = append(walker.filters, toFilter(f))
			return nil
		})
		flag.Func("iregex", "", func(s string) error {
			f, err := filter.NewIRegex(s)
			if err != nil {
				return err
			}
			walker.filters = append(walker.filters, toFilter(f))
			return nil
		})
		flag.Func("type", "", func(s string) error {
			f, err := filter.NewFileType(s)
			if err != nil {
				return err
			}
			walker.filters = append(walker.filters, toFilter(f))
			return nil
		})
		flag.Var(boolFunc(func(b bool) {
			if b {
				walker.prunes = append(walker.prunes, walker.filters...)
				walker.filters = walker.filters[:0]
			}
		}), "prune", "")
	}

	roots, remain := getRoots(args[1:])
	if err := flag.Parse(remain); err != nil {
		return nil, nil, err
	}
	if !isGC {
		debug.SetGCPercent(-1)
	}
	return walker, roots, nil
}

func toFilter(f filter.FileExp) filter.FileExp {
	if isNot {
		isNot = false
		return filter.NewNotFilter(f)
	}
	if isNotProcess {
		return filter.NewNotFilter(f)
	}
	return f
}

func getRoots(args []string) (roots []string, remain []string) {
	roots = make([]string, 0, defaultMakeLen)
	for i, arg := range args {
		if len(arg) == 0 {
			break
		}
		if arg[0] == '-' {
			break
		}
		roots = append(roots, arg)
		remain = args[i+1:]
	}

	if len(roots) == 0 {
		return []string{"."}, args
	}

	return roots, remain
}
