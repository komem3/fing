# Fing

Fing is a fast file finder.
It doesn't cover all of find's rich options, but fing provides a similar interface.

## Why fing

The `find` command supports a rich of useful options, but it is very slow.
On the other hand, [`fd`](https://github.com/sharkdp/fd) can search very fast, but it is very different from how to use find.

So I created fing with the goal of creating a fast `find` command with an interface similar to `find`.
If you want a fast `find` command, fing is a good choice.

## Install

### Go Tool

```bash
go install github.com/komem3/fing@latest
```

### From Source Code

```bash
# in project directory
go build -o fing .

sudo mv fing /usr/bin/
```

### Usage

```
Usage: fing [staring-point...] [flag] [expression]

Fing is a fast file finder that provides an interface similar to find.

flags are:
  -dry
    Only output parse result of expression.
    If this option is specified, the file will not be searched.
  -maxdepth
    The depth to search.
    Unlike find, it can be specified at the same time as prune.
  -EI
    Exclude pattern from I option.
    This uses the before expressions as well as prune.
    example: -I <expression> -EI
  -I
    Ignore files in .gitignore.

expression are:
  -a -and
    This flag is skipped.
  -empty
    Search emptry file and directory.
    This is shothand of '-size 0c'.
  -executable
    Match files which are executable by current user.
  -false
    Always false.
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
  -print
    Add a new line character after the file name. This option is default enabled.
  -print0
    Add a null character after the file name.
  -prune
    Prunes directory that match before expressions.
    example: <expression> -prune
  -regex string
    Search for files using regular expressions.
    This option match to file path.
  -rname string
    Search for files using regular expressions.
    This option match only to file name.
  -size [+|-]n[ckMG]
    The size of file. Should specify the unit of size.
    c(for bytes), k(for KiB), M(for MiB), G(for Gib).
  -true
    Always true.
  -type string
    File is type.
    Support file(f), directory(d), named piep(p) and socket(s).
```

### Examples

- search jpg files.

```bash
fing ./testdata -name "*.jpg"
```

- ignores files in .gitignore and hidden files (search like `fd`).

```bash
# fd txt ./testdata
fing ./testdata -I -name ".*" -prune -false -o -not -name ".*" -irname ".*txt.*"
```

- Each operator is AND expression, but you can also specify OR expression.

```bash
fing ./testdata -name "*.jpg" -o -name "*.png"
```

- Debug option `-dry`. You can see how `fing` evaluated the expression.

```bash
fing -dry -name "*.jpg" -name "*.png"
```

## NOTE

- The regular expression uses Go's [regexp](https://pkg.go.dev/regexp) package, so it behaves differently than the find command's regular expression.
- The find command strictly evaluates operators from left to right, but fing may not do so to optimize process. Therefore, the following options may behave differently than find.
  - print
  - print0
  - prune

## Benchmark

```
hyperfine -i "find . -iname '*[0-9].jpg'" "fdfind -HI '.*[0-9]\.jpg$'" "fing . -iname '*[0-9].jpg'"
```

| Command                      |      Mean [s] | Min [s] | Max [s] |    Relative |
| :--------------------------- | ------------: | ------: | ------: | ----------: |
| `find . -iname '*[0-9].jpg'` | 8.405 ± 4.451 |   5.180 |  15.911 | 6.09 ± 3.25 |
| `fdfind -HI '.*[0-9]\.jpg$'` | 1.379 ± 0.081 |   1.315 |   1.579 |        1.00 |
| `fing . -iname '*[0-9].jpg'` | 1.432 ± 0.069 |   1.346 |   1.604 | 1.04 ± 0.08 |

## Author

komem3

## License

MIT
