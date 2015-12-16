package integration

import (
	"errors"
	"path"

	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

type MapFileReader struct {
	files map[string][]byte
}

func NewMapFileReader(files map[string][]byte) MapFileReader {
	return MapFileReader{files}
}

func (m MapFileReader) ReadFileContent(file string) ([]byte, error) {
	if content, ok := m.files[path.Base(file)]; ok {
		return content, nil
	}
	return nil, errors.New("can't find file")
}

func (m MapFileReader) ReadPath(path string) ([]*file.File, error) {
	result := make([]*file.File, 0, len(m.files))
	for k := range m.files {
		if file, err := file.NewFile(k, ""); err == nil {
			result = append(result, file)
		}
	}

	return result, nil
}
