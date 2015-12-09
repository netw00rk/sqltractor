package integration

import (
	"os"
)

type MockFileInfo struct {
	Name string
}

func (fi MockFileInfo) Name() string {
	return fi.Name
}

func (fi MockFileInfo) Size() int64 {
	return 0
}

func (fi MockFileInfo) Mode() os.FileMode {
	return os.ModeAppend
}

func (fi MockFileInfo) ModeTime()

type MapFileReader struct {
	files map[string]byte
}

func NewMapFileReader(files map[string]byte) MapFileReader {
	return MapFileReader(files)
}

func (m MockFileReader) ReadFileContent(path string) ([]byte, error) {
	return []byte("test"), nil
}

func (m MockFileReader) ReadPath(path string) ([]os.FileInfo, error) {
	return nil, nil
}
