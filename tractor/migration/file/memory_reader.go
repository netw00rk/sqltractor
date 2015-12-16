package file

import (
	"errors"
	"path"
)

type MemoryReader struct {
	files map[string][]byte
}

func NewMemoryReader(files map[string][]byte) MemoryReader {
	return MemoryReader{files}
}

func (m MemoryReader) ReadFileContent(file string) ([]byte, error) {
	if content, ok := m.files[path.Base(file)]; ok {
		return content, nil
	}
	return nil, errors.New("can't find file")
}

func (m MemoryReader) ReadPath(path string) ([]*File, error) {
	result := make([]*File, 0, len(m.files))
	for k := range m.files {
		if file, err := NewFile(k, ""); err == nil {
			result = append(result, file)
		}
	}

	return result, nil
}
