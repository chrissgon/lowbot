package lowbot

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

var (
	FILE_AUDIO_MOCK    = NewFile("./mocks/music.mp3")
	FILE_DOCUMENT_MOCK = NewFile("./mocks/features.txt")
	FILE_IMAGE_MOCK    = NewFile("./mocks/image.jpg")
	FILE_VIDEO_MOCK    = NewFile("./mocks/video.mp4")
)

func TestFile_NewFile(t *testing.T) {
	path := "./mocks/features.txt"
	p, err := filepath.Abs(path)
	fileID := uuid.New()

	expect := &File{
		FileID:    fileID,
		FileType:  FILETYPE_DOCUMENT,
		Extension: ".txt",
		Path:      p,
		Err:       err,
	}

	have := NewFile(path)
	have.GetFile().FileID = fileID

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFile_GetFile(t *testing.T) {
	expect := NewFile("./mocks/features.txt")
	have := expect.GetFile()

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFile_Read(t *testing.T) {
	file := NewFile("./mocks/features.txt")

	err := file.Read()

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if file.GetFile().Bytes == nil {
		t.Errorf(FormatTestError([]byte{}, nil))
	}

	file = NewFile("https://myurl.com/features.txt")

	err = file.Read()

	if !errors.Is(err, ERR_FEATURE_UNIMPLEMENTED) {
		t.Errorf(FormatTestError(ERR_FEATURE_UNIMPLEMENTED, err))
	}
}

func TestFile_SetFileType(t *testing.T) {
	expect := FILETYPE_AUDIO
	have := FILE_AUDIO_MOCK.GetFile().FileType

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	expect = FILETYPE_DOCUMENT
	have = FILE_DOCUMENT_MOCK.GetFile().FileType

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	expect = FILETYPE_IMAGE
	have = FILE_IMAGE_MOCK.GetFile().FileType

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	expect = FILETYPE_VIDEO
	have = FILE_VIDEO_MOCK.GetFile().FileType

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFile_IsAudio(t *testing.T) {
	file := &File{
		FileType: FILETYPE_AUDIO,
	}

	expect := true
	have := file.IsAudio()

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFile_IsDocument(t *testing.T) {
	file := &File{
		FileType: FILETYPE_DOCUMENT,
	}

	expect := true
	have := file.IsDocument()

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFile_IsImage(t *testing.T) {
	file := &File{
		FileType: FILETYPE_IMAGE,
	}

	expect := true
	have := file.IsImage()

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFile_IsVideo(t *testing.T) {
	file := &File{
		FileType: FILETYPE_VIDEO,
	}

	expect := true
	have := file.IsVideo()

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}
