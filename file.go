package lowbot

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type IFile interface {
	GetFile() *File
	Read() error
}

type File struct {
	FileID    uuid.UUID
	Bytes     []byte
	Path      string
	Extension string
	Err       error
}

func NewFile(path string) IFile {
	file := &File{
		FileID: uuid.New(),
	}

	file.Extension = filepath.Ext(path)
	file.Path = path

	return file
}

func (file *File) Read() error {
	if IsURL(file.Path) {
		file.Err = errors.New("file Read url unimplemented")
		return file.Err
	}

	file.Bytes, file.Err = os.ReadFile(file.Path)

	return file.Err
}

func (file *File) GetFile() *File {
	return file
}
