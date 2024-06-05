package lowbot

import (
	"path/filepath"
	"testing"
)

func TestFile_NewFile(t *testing.T) {
	path := "./mocks/features.txt"
	file := NewFile(path)

	expect := ".txt"
	have := file.GetFile().Extension

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	expect, _ = filepath.Abs(path)
	have = file.GetFile().Path

	if expect != have {
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
}
