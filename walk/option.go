package walk

import (
	"flag"
	"fmt"
	"io"
	"os"
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

Fing is A fast file finder that provides an interface similar to find.

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
  -iname string
    Like -name, but the match is case insensitive.
  -ipath string
    Like -path, but the match is case insensitive.
  -iregex string
    Like -regex, but the match is case insensitive.
  -irname string
    Like -rname, but the match is case insensitive.
  -name string
    Search for files using wildcard expressions.
    This option match only to file name.
  -not
    True if next expression false.
  -or
  -o
    Evaluate the previous and next expressions with or.
  -path string
    Search for files using wildcard expressions.
    This option match to file path.
  -prune
    Prunes directory that match before expressions.
  -regex string
    Search for files using regular expressions.
    This option match to file path.
    Unlike find, this is a backward match.
  -rname string
    Search for files using regular expressions.
    This option match only to file name..
    Unlike regex option, this option is exact match.
  -type string
    File is type.
    Support file(f), directory(d), named piep(p) and socket(s).
`

func NewWalkerFromArgs(args []string, out, outerr io.Writer) (*Walker, directoryInfos, error) {
	walker := &Walker{
		matcher:     make(filter.OrExp, 0, defaultMakeLen),
		prunes:      make(filter.OrExp, 0, defaultMakeLen),
		out:         out,
		outerr:      outerr,
		depth:       -1,
		targets:     make(directoryInfos, 0, defaultDirecotryBuffer),
		concurrency: make(chan struct{}, concurrencyMax),
	}

	flag := flag.NewFlagSet(args[0], flag.ExitOnError)
	flag.Usage = func() { fmt.Fprint(os.Stderr, Usage) }

	exp := make(filter.AndExp, 0, defaultMakeLen)
	{
		var isNot bool
		// expression
		flag.BoolVar(&isNot, "not", false, "")
		flag.Func("name", "", func(s string) error {
			exp = append(exp, toFilter(filter.NewFileName(s), &isNot))
			return nil
		})
		flag.Func("iname", "", func(s string) error {
			exp = append(exp, toFilter(filter.NewIFileName(s), &isNot))
			return nil
		})
		flag.Func("path", "", func(s string) error {
			exp = append(exp, toFilter(filter.NewPath(s), &isNot))
			return nil
		})
		flag.Func("ipath", "", func(s string) error {
			exp = append(exp, toFilter(filter.NewIPath(s), &isNot))
			return nil
		})
		flag.Func("regex", "", func(s string) error {
			f, err := filter.NewRegex(s)
			if err != nil {
				return err
			}
			exp = append(exp, toFilter(f, &isNot))
			return nil
		})
		flag.Func("iregex", "", func(s string) error {
			f, err := filter.NewIRegex(s)
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
		flag.Func("irname", "", func(s string) error {
			f, err := filter.NewIRegexName(s)
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
		flag.Var(boolFunc(func(b bool) {
			if b {
				walker.matcher = append(walker.matcher, exp)
				walker.prunes = append(walker.prunes, walker.matcher...)
				exp = make(filter.AndExp, 0, defaultMakeLen)
				walker.matcher = walker.matcher[:0]
			}
		}), "prune", "")
		flag.Var(boolFunc(func(b bool) {
			if b {
				walker.matcher = append(walker.matcher, exp)
				exp = make(filter.AndExp, 0, defaultMakeLen)
			}
		}), "or", "")
		flag.Var(boolFunc(func(b bool) {
			if b {
				walker.matcher = append(walker.matcher, exp)
				exp = make(filter.AndExp, 0, defaultMakeLen)
			}
		}), "o", "")
	}

	roots, remain := getRoots(args[1:])
	remain, err := setOption(walker, remain)
	if err != nil {
		return nil, nil, err
	}
	if err := flag.Parse(remain); err != nil {
		return nil, nil, err
	}
	walker.matcher = append(walker.matcher, exp)
	return walker, roots, nil
}

func toFilter(f filter.FileExp, isNot *bool) filter.FileExp {
	if *isNot {
		*isNot = false
		return filter.NewNotExp(f)
	}
	return f
}

func setOption(walker *Walker, args []string) (remine []string, err error) {
	remine = args[:]
	for len(remine) > 0 && len(remine[0]) > 0 {
		switch remine[0] {
		case "-I":
			walker.gitignore = true
		case "-dry":
			walker.IsDry = true
		case "-maxdepth":
			if len(remine) < 2 {
				return remine, nil
			}
			d, err := strconv.Atoi(remine[1])
			if err != nil {
				return nil, err
			}
			walker.depth = d
			remine = remine[1:]
		default:
			return remine, nil
		}
		if len(remine) == 1 {
			return nil, nil
		}
		remine = remine[1:]
	}
	return remine, nil
}

func getRoots(args []string) (roots []*direcotryInfo, remain []string) {
	roots = make(directoryInfos, 0, defaultMakeLen)
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

	if len(roots) == 0 {
		return directoryInfos{{path: "."}}, args
	}

	return roots, remain
}
