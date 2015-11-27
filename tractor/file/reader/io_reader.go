package reader

import (
	"io/ioutil"
	"os"
)


type IOFileReader struct {
}

func (IOFileReader) ReadFileContent(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (IOFileReader) ReadPath(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}
