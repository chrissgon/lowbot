package lowbot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

var CommonMIMETypes = map[string]string{
	// Images
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".tiff": "image/tiff",
	".ico":  "image/x-icon",
	".svg":  "image/svg+xml",
	".apng": "image/apng",
	".avif": "image/avif",

	// Audio
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".ogg":  "audio/ogg",
	".oga":  "audio/ogg",
	".opus": "audio/opus",
	".aac":  "audio/aac",
	".m4a":  "audio/mp4",
	".weba": "audio/webm",
	".flac": "audio/flac",

	// Video
	".mp4":  "video/mp4",
	".mov":  "video/quicktime",
	".wmv":  "video/x-ms-wmv",
	".avi":  "video/x-msvideo",
	".webm": "video/webm",
	".flv":  "video/x-flv",
	".mkv":  "video/x-matroska",
	".mpeg": "video/mpeg",
	".ogv":  "video/ogg",

	// Documents
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".odt":  "application/vnd.oasis.opendocument.text",

	// Text & Code
	".txt":  "text/plain",
	".csv":  "text/csv",
	".html": "text/html",
	".htm":  "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".json": "application/json",
	".xml":  "application/xml",
	".yaml": "text/yaml",
	".yml":  "text/yaml",

	// Archives
	".zip": "application/zip",
	".rar": "application/vnd.rar",
	".tar": "application/x-tar",
	".gz":  "application/gzip",
	".7z":  "application/x-7z-compressed",

	// Fonts
	".ttf":   "font/ttf",
	".otf":   "font/otf",
	".woff":  "font/woff",
	".woff2": "font/woff2",

	// Others
	".apk": "application/vnd.android.package-archive",
	".exe": "application/vnd.microsoft.portable-executable",
	".swf": "application/x-shockwave-flash",
}

var (
	FILETYPE_AUDIO_EXT = []string{".mp3", ".wav", ".ogg", ".oga", ".opus", ".aac", ".m4a", ".weba", ".flac"}
	FILETYPE_IMAGE_EXT = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".ico", ".svg", ".apng", ".avif"}
	FILETYPE_VIDEO_EXT = []string{".mp4", ".mov", ".wmv", ".avi", ".webm", ".flv", ".mkv", ".mpeg", ".ogv"}
)

func NewFile(path, url string) IFile {
	file := &File{
		FileID:   uuid.New(),
		FileType: FILETYPE_DOCUMENT,
		URL:      url,
	}

	file.Extension = filepath.Ext(path)
	file.Name = filepath.Base(path)
	file.Mime = GetMIMEByExtension(path)

	file.SetFilePath(path)
	file.SetFileType()

	return file
}

func (file *File) GetFile() *File {
	return file
}

func (file *File) Read() error {
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
		return fmt.Errorf("ERROR TO GET FILE: %s: %s", resp.Status, file.URL)
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

func GetMIMEByExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if mime, ok := CommonMIMETypes[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}
