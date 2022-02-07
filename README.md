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
```

### Examples

- search jpg files.

```bash
fing ./testdata -name "*.jpg"
```

- ignores files in .gitignore and hidden files (search like `fd`).

```bash
# fd txt ./testdata
fing ./testdata -I -name ".*" -prune -not -name ".*" -irname ".*txt.*"
```

- Each operator is AND expression, but you can also specify OR expression.

```bash
fing ./testdata -name "*.jpg" -o -name "*.png"
```

- Debug option `-dry`. You can see how `fing` evaluated the expression.

```bash
fing -dry -name "*.jpg" -name "*.png"
```

## Benchmark

| Command                      |      Mean [s] | Min [s] | Max [s] |    Relative |
| :--------------------------- | ------------: | ------: | ------: | ----------: |
| `find . -iname '*[0-9].jpg'` | 4.181 ± 0.085 |   4.091 |   4.336 | 3.54 ± 0.25 |
| `fdfind -HI '.*[0-9]\.jpg$'` | 1.192 ± 0.012 |   1.173 |   1.216 | 1.01 ± 0.07 |
| `fing . -iname '*[0-9].jpg'` | 1.180 ± 0.079 |   1.122 |   1.377 |        1.00 |

## Author

komem3

## License

MIT
