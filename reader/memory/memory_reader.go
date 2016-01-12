package memory

import (
	"errors"

	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

type MemoryReader struct {
	files map[string][]byte
}

func NewMemoryReader(files map[string][]byte) *MemoryReader {
	return &MemoryReader{files}
}

func (r *MemoryReader) Read() ([]*file.File, error) {
	result := make([]*file.File, 0, len(r.files))
	for k := range r.files {
		if file, err := file.NewFile(k, r.buildContentFunc(k)); err == nil {
			result = append(result, file)
		}
	}

	return result, nil
}

func (r *MemoryReader) buildContentFunc(name string) func() ([]byte, error) {
	return func() ([]byte, error) {
		if content, ok := r.files[name]; ok {
			return content, nil
		}
		return nil, errors.New("can't find file")
	}
}
