package lowbot

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"slices"

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
	Name      string
	Bytes     []byte
	Path      string
	URL       string
	Extension string
	Mime      string
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

func NewFile(path, url string) IFile {
	fmt.Println("test", path, url)
	file := &File{
		FileID:   uuid.New(),
		FileType: FILETYPE_DOCUMENT,
		URL:      url,
	}

	file.Extension = filepath.Ext(path)
	file.Name = filepath.Base(path)
	file.Mime = mime.TypeByExtension(file.Extension)

	file.SetFilePath(path)
	file.SetFileType()

	return file
}

func (file *File) GetFile() *File {
	return file
}

func (file *File) Read() error {
	fmt.Println("url", file.URL)
	fmt.Println("path", file.Path)
	if file.URL == "" && !IsURL(file.URL) {
		file.Bytes, file.Err = os.ReadFile(file.Path)
		return file.Err
	}

	resp, err := http.Get(file.URL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ERROR TO GET FILE: %s", resp.Status)
	}

	file.Bytes, file.Err = io.ReadAll(resp.Body)

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
	if slices.Contains(FILETYPE_AUDIO_EXT, file.Extension) {
		file.FileType = FILETYPE_AUDIO
		return
	}
	if slices.Contains(FILETYPE_IMAGE_EXT, file.Extension) {
		file.FileType = FILETYPE_IMAGE
		return
	}
	if slices.Contains(FILETYPE_VIDEO_EXT, file.Extension) {
		file.FileType = FILETYPE_VIDEO
		return
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
