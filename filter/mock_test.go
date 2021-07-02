package filter

import "io/fs"

type mockDirEntry struct {
	name  string
	isDir bool
	typ   fs.FileMode
	info  fs.FileInfo
}

func NewMockDriEntry(name string, isdir bool, typ fs.FileMode, info fs.FileInfo) fs.DirEntry {
	return &mockDirEntry{name, isdir, typ, info}
}

var _ fs.DirEntry = (*mockDirEntry)(nil)

func (m *mockDirEntry) Name() string {
	return m.name
}

func (m *mockDirEntry) IsDir() bool {
	return m.isDir
}

func (m *mockDirEntry) Type() fs.FileMode {
	return m.typ
}

func (m *mockDirEntry) Info() (fs.FileInfo, error) {
	return m.info, nil
}
