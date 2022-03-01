package walk

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/komem3/fing/filter"
)

const defaultMakeLen = 1 << 2

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
Usage: fing [staring-point...] [flag] [expression]

Fing is a fast file finder that provides an interface similar to find.

flags are:
  -I
    Ignore files in .gitignore.
  -dry
    Only output parse result of expression.
    If this option is specified, the file will not be searched.
  -maxdepth
    The depth to search.
    Unlike find, it can be specified at the same time as prune.

expression are:
  -a -and
    This flag is skipped.
  -empty
    Search emptry file and directory.
    This is shothand of '-size 0c'.
  -executable
    Match files which are executable by current user.
  -iname string
    Like -name, but the match is case insensitive.
  -ipath string
    Like -path, but the match is case insensitive.
  -iregex string
    Like -regex, but the match is case insensitive.
  -irname string
    Like -rname, but the match is case insensitive.
  -name string
    Search for files using glob expressions.
    This option match only to file name.
  -not
    True if next expression false.
  -o -or
    Evaluate the previous and next expressions with or.
  -path string
    Search for files using wildcard expressions.
    This option match to file path.
    Unlike find, This option explicitly matched by using one or more <slash>.
  -print0
    Add a null character after the file name.
  -prune
    Prunes directory that match before expressions.
  -regex string
    Search for files using regular expressions.
    This option match to file path.
  -rname string
    Search for files using regular expressions.
    This option match only to file name.
  -size [+|-]n[ckMG]
    The size of file. Should specify the unit of size.
    c(for bytes), k(for KiB), M(for MiB), G(for Gib).
  -type string
    File is type.
    Support file(f), directory(d), named piep(p) and socket(s).
`

func NewWalkerFromArgs(args []string, out, outerr *bufio.Writer) (*Walker, directoryInfos, error) {
	walker := &Walker{
		matcher:     make(filter.OrExp, 0, defaultMakeLen),
		prunes:      make(filter.OrExp, 0, defaultMakeLen),
		out:         out,
		outerr:      outerr,
		depth:       -1,
		targets:     make(directoryInfos, 0, defaultDirecotryBuffer),
		concurrency: make(chan struct{}, concurrencyMax),
		printType:   println,
	}

	flag := flag.NewFlagSet(args[0], flag.ExitOnError)
	flag.Usage = func() { fmt.Fprint(os.Stderr, Usage) }

	{
		// flags
		flag.BoolVar(&walker.gitignore, "I", false, "")
		flag.BoolVar(&walker.IsDry, "dry", false, "")
		flag.Func("maxdepth", "", func(s string) error {
			d, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			walker.depth = d
			return nil
		})
	}

	exp := make(filter.AndExp, 0, defaultMakeLen)
	{
		// expression
		var isNot bool
		_ = flag.Bool("a", false, "")
		_ = flag.Bool("and", false, "")
		flag.Var(boolFunc(func(b bool) {
			if b {
				f, err := filter.NewSize("0c")
				if err != nil {
					panic(err)
				}
				exp = append(exp, toFilter(f, &isNot))
			}
		}), "empty", "")
		flag.Var(boolFunc(func(b bool) {
			if b {
				exp = append(exp, toFilter(filter.NewExecutable(), &isNot))
			}
		}), "executable", "")
		flag.Func("iname", "", func(s string) error {
			f, err := filter.NewIFileName(s)
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Func("ipath", "", func(s string) error {
			f, err := filter.NewIPath(filepath.FromSlash(s))
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Func("iregex", "", func(s string) error {
			f, err := filter.NewIRegex(filepath.FromSlash(s))
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Func("irname", "", func(s string) error {
			f, err := filter.NewIRegexName(s)
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Func("name", "", func(s string) error {
			f, err := filter.NewFileName(s)
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.BoolVar(&isNot, "not", false, "")
		orFunc := func(b bool) {
			if b {
				walker.matcher = append(walker.matcher, exp)
				exp = make(filter.AndExp, 0, defaultMakeLen)
			}
		}
		flag.Var(boolFunc(orFunc), "o", "")
		flag.Var(boolFunc(orFunc), "or", "")
		flag.Func("path", "", func(s string) error {
			f, err := filter.NewPath(filepath.FromSlash(s))
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Var(boolFunc(func(b bool) {
			if b {
				walker.printType = print0
			}
		}), "print0", "")
		flag.Var(boolFunc(func(b bool) {
			if b {
				walker.matcher = append(walker.matcher, exp)
				walker.prunes = append(walker.prunes, walker.matcher...)
				exp = make(filter.AndExp, 0, defaultMakeLen)
				walker.matcher = walker.matcher[:0]
			}
		}), "prune", "")
		flag.Func("regex", "", func(s string) error {
			f, err := filter.NewRegex(filepath.FromSlash(s))
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Func("rname", "", func(s string) error {
			f, err := filter.NewRegexName(s)
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Func("size", "", func(s string) error {
			f, err := filter.NewSize(s)
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Func("type", "", func(s string) error {
			f, err := filter.NewFileType(s)
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
	}

	roots, remain := getRoots(args[1:], false)
	if err := flag.Parse(remain); err != nil {
		return nil, nil, err
	}
	backRoots, _ := getRoots(flag.Args(), len(roots) == 0)

	walker.matcher = append(walker.matcher, exp)
	return walker, append(roots, backRoots...), nil
}

func toFilter(f filter.FileExp, isNot *bool) filter.FileExp {
	if *isNot {
		*isNot = false
		return filter.NewNotExp(f)
	}
	return f
}

func getRoots(args []string, leastOne bool) (roots []*direcotryInfo, remain []string) {
	remain = args[:]
	for i, arg := range args {
		if len(arg) == 0 {
			break
		}
		if arg[0] == '-' {
			break
		}
		roots = append(roots, &direcotryInfo{path: arg})
		remain = args[i+1:]
	}

	if leastOne && len(roots) == 0 {
		return directoryInfos{{path: "."}}, args
	}

	return roots, remain
}
