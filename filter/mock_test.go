package filter_test

import (
	"io/fs"
	"time"
)

type mockDirFileInfo struct {
	name    string
	isDir   bool
	typ     fs.FileMode
	size    int64
	modTime time.Time
	sys     interface{}
}

var (
	_ fs.DirEntry = (*mockDirFileInfo)(nil)
	_ fs.FileInfo = (*mockDirFileInfo)(nil)
)

func (m *mockDirFileInfo) Name() string {
	return m.name
}

func (m *mockDirFileInfo) IsDir() bool {
	return m.isDir
}

func (m *mockDirFileInfo) Type() fs.FileMode {
	return m.typ
}

func (m *mockDirFileInfo) Mode() fs.FileMode {
	return m.typ
}

func (m *mockDirFileInfo) Info() (fs.FileInfo, error) {
	return m, nil
}

func (m *mockDirFileInfo) Size() int64 {
	return m.size
}

func (m *mockDirFileInfo) ModTime() time.Time {
	return m.modTime
}

func (m *mockDirFileInfo) Sys() interface{} {
	return m.sys
}
