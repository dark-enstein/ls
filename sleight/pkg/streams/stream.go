package streams

import (
	"bytes"
	"os"
)

type Stream interface {
	Store(name string, b []byte) error
	Retrieve(name string) (int, error)
}

type File struct {
	Name string
	Buffer
	os.File
}

type Buffer struct {
	b bytes.Buffer
}


func NewFile(path string) (*File, error) {
	if isNotExist
	f := &File{
		Name:   path,
		Buffer: Buffer{},
		File:   os.File{},
	}
}