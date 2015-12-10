package file

type Reader interface {
	ReadFileContent(path string) ([]byte, error)
	ReadPath(path string) ([]*File, error)
}
