package file

import (
	"io/ioutil"
	"path"

	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

type FileReader struct {
	Path string
}

func NewFileReader(path string) *FileReader {
	return &FileReader{path}
}

func (r *FileReader) Read() ([]*file.File, error) {
	ioFiles, err := ioutil.ReadDir(r.Path)
	if err != nil {
		return nil, err
	}

	files := make([]*file.File, 0)
	for _, ioFile := range ioFiles {
		if file, err := file.NewFile(ioFile.Name(), r.buildContentFunc(ioFile.Name())); err == nil {
			files = append(files, file)
		}
	}

	return files, nil
}

func (r *FileReader) buildContentFunc(name string) func() ([]byte, error) {
	return func() ([]byte, error) {
		return ioutil.ReadFile(path.Join(r.Path, name))
	}
}
