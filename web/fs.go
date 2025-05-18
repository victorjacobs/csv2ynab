package web

import (
	"context"
	"io/fs"

	"golang.org/x/net/webdav"
)

type FileHandler func(path string, contents []byte)

type FS struct {
	webdav.FileSystem

	handleFile FileHandler
}

var _ webdav.FileSystem = &FS{}

func NewFS(handleFile FileHandler) *FS {
	return &FS{
		FileSystem: webdav.NewMemFS(),
		handleFile: handleFile,
	}
}

func (f *FS) Mkdir(ctx context.Context, name string, perm fs.FileMode) error {
	return nil
}

func (f *FS) OpenFile(ctx context.Context, name string, flag int, perm fs.FileMode) (webdav.File, error) {
	file, err := f.FileSystem.OpenFile(ctx, name, flag, perm)
	if err != nil {
		return nil, err
	}

	return &File{
		File: file,
		path: name,
		handleFile: func(path string, contents []byte) {
			f.handleFile(path, contents)

			f.RemoveAll(ctx, name)
		},
	}, nil
}

func (f *FS) Rename(ctx context.Context, oldName string, newName string) error {
	return nil
}
