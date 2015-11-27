package reader

import "os"

type FileReader interface {
	ReadFileContent(path string) ([]byte, error)
	ReadPath(path string) ([]os.FileInfo, error)
}