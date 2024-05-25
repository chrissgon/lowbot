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
	FileType  FileType
	Bytes     []byte
	Path      string
	Extension string
	Err       error
}

type FileType string

const (
	FILETYPE_AUDIO    FileType = "audio"
	FILETYPE_DOCUMENT FileType = "document"
	FILETYPE_IMAGE    FileType = "image"
	FILETYPE_VIDEO    FileType = "video"
)

var (
	FILETYPE_AUDIO_EXT = []string{".aac", ".mp3", ".oga", ".opus", ".wav", ".weba", ".cda"}
	FILETYPE_IMAGE_EXT = []string{".apng", ".avif", ".gif", ".jpg", ".jpeg", ".png", ".svg", ".webp"}
	FILETYPE_VIDEO_EXT = []string{".avi", ".mp4", ".mpeg", ".ogv", ".webm"}
)

func NewFile(path string) IFile {
	var err error

	file := &File{
		FileID:   uuid.New(),
		FileType: FILETYPE_DOCUMENT,
	}

	file.Extension = filepath.Ext(path)
	file.Path, err = filepath.Abs(path)

	if err != nil {
		file.Path = path
	}

	file.SetFileType()

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

func (file *File) SetFileType() {
	for _, ext := range FILETYPE_AUDIO_EXT {
		if file.Extension == ext {
			file.FileType = FILETYPE_AUDIO
			return
		}
	}
	for _, ext := range FILETYPE_IMAGE_EXT {
		if file.Extension == ext {
			file.FileType = FILETYPE_IMAGE
			return
		}
	}
	for _, ext := range FILETYPE_VIDEO_EXT {
		if file.Extension == ext {
			file.FileType = FILETYPE_VIDEO
			return
		}
	}
}
