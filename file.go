package lowbot

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type IFile interface {
	GetFile() *File
	Read() error
	IsAudio() bool
	IsDocument() bool
	IsImage() bool
	IsVideo() bool
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
	FILETYPE_AUDIO    FileType = "FILETYPE_AUDIO"
	FILETYPE_DOCUMENT FileType = "FILETYPE_DOCUMENT"
	FILETYPE_IMAGE    FileType = "FILETYPE_IMAGE"
	FILETYPE_VIDEO    FileType = "FILETYPE_VIDEO"
)

var (
	FILETYPE_AUDIO_EXT = []string{".aac", ".mp3", ".oga", ".opus", ".wav", ".weba", ".cda"}
	FILETYPE_IMAGE_EXT = []string{".apng", ".avif", ".gif", ".jpg", ".jpeg", ".png", ".svg", ".webp"}
	FILETYPE_VIDEO_EXT = []string{".avi", ".mp4", ".mpeg", ".ogv", ".webm"}
)

func NewFile(path string) IFile {
	file := &File{
		FileID:   uuid.New(),
		FileType: FILETYPE_DOCUMENT,
	}

	file.Extension = filepath.Ext(path)

	file.SetFilePath(path)
	file.SetFileType()

	return file
}

func (file *File) GetFile() *File {
	return file
}

func (file *File) Read() error {
	if IsURL(file.Path) {
		file.Err = ERR_FEATURE_UNIMPLEMENTED
		return file.Err
	}

	file.Bytes, file.Err = os.ReadFile(file.Path)

	return file.Err
}

func (file *File) SetFilePath(path string) {
	if IsURL(path) {
		file.Path = path
		return
	}

	file.Path, file.Err = filepath.Abs(path)

	if file.Err != nil {
		file.Path = path
	}
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

func (file *File) IsAudio() bool {
	return file.FileType == FILETYPE_AUDIO
}

func (file *File) IsDocument() bool {
	return file.FileType == FILETYPE_DOCUMENT
}

func (file *File) IsImage() bool {
	return file.FileType == FILETYPE_IMAGE
}

func (file *File) IsVideo() bool {
	return file.FileType == FILETYPE_VIDEO
}
