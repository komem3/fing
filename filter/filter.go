package filter

import "io/fs"

type FileExp interface {
	Match(path string, info fs.DirEntry) (bool, error)
}
