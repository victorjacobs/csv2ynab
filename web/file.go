package web

import (
	"io"

	"golang.org/x/net/webdav"
)

type File struct {
	webdav.File

	path       string
	handleFile FileHandler
}

func (f *File) Close() error {
	if err := f.File.Close(); err != nil {
		return err
	}

	// Rewind file
	f.Seek(0, 0)

	if fileContents, err := io.ReadAll(f.File); err == nil {
		f.handleFile(f.path, fileContents)
	}

	return nil
}
