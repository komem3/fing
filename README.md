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
hyperfine --warmup 3  'fing -type f -iname "*.jpg"' 'find -type f -iname "*.jpg"' 'fdfind -I -H -g -t f "*.jpg"'
```

## TODO

- not option の実装
- type option の実装
- perm option の実装
- exec の実装
