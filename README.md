# Fing

Fing is A fast file finder.
It doesn't cover all of find's rich options, but fing provides a similar interface.

## Why fing

The `find` command supports a rich of useful options, but it is very slow.
On the other hand, [`fd`](https://github.com/sharkdp/fd) can search very fast, but it is very different from how to use find.

So I created fing with the goal of creating a fast `find` command with an interface similar to `find`.
If you want a fast `find` command, fing is a good choice.

## Install

### Go Tool

```bash
go install github.com/komem3/fing
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

Fing is A fast file finder that provides an interface similar to find.

flags are:
  -I
    Ignore files in .gitignore.
    This is a fing specific option.
  -dry
    Only output parse result of expression.
    If this option is specified, the file will not be searched.
    This is a fing specific option.
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
    This is a fing specific option.
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
    This is a fing specific option.
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

- ignores files in .gitignore and hidden files (search like fd).

```bash
# fd txt ./testdata
fing ./testdata -I -name ".*" -prune -not -name ".*" -irname ".*txt.*"
```

- Each operator is AND expression, but you can also specify OR expression.

```bash
fing ./testdata -name "*.jpg" -or -name "*.png"
```

- Debug option `-dry`. You can see how fing evaluated the expression.

```bash
fing -dry -name "*.jpg" -name "*.png"
```

## Author

komem3

## Licence

MIT
