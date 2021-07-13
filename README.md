# Fing

Fing is a simple and fast file finder.
It doesn't cover all of find's rich options, but aim to use it in the same way.

## Build

```shell
# in project directory
go build -o fing .

sudo mv fing /usr/bin/
```

## Benchmark

- fdfind vs find vs fing

```shell
hyperfine --warmup 3 'find -type f -iname "*.jpg"' 'fdfind -I -H -g -t f "*.jpg"' 'fing -type f -iname "*.jpg"'
```

## TODO

- exec の実装
