package file

import "io/ioutil"

type IOReader struct {
}

func (IOReader) ReadFileContent(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (IOReader) ReadPath(path string) ([]*File, error) {
	ioFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := make([]*File, 0)
	for _, ioFile := range ioFiles {
		if file, err := NewFile(ioFile.Name(), path); err == nil {
			files = append(files, file)
		}
	}

	return files, nil
}
