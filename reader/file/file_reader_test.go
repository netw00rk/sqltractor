package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FileReaderTestSuite struct {
	suite.Suite
	path string
}

func (s *FileReaderTestSuite) SetupSuite() {
	s.path, _ = ioutil.TempDir("/tmp", "TestLookForMigrationFilesInSearchPath")
	ioutil.WriteFile(path.Join(s.path, "001_migrationfile.up.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "001_migrationfile.down.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "002_migrationfile.up.sql"), []byte("test"), 0755)
}

func (s *FileReaderTestSuite) TearDownSuite() {
	os.RemoveAll(s.path)
}

func (s *FileReaderTestSuite) TestReadFiles() {
	reader := NewFileReader(s.path)
	files, err := reader.Read()
	s.Equal(3, len(files))

	content, err := files[2].Content()
	s.Nil(err, fmt.Sprintf("Can't read %s file content", files[2].FileName))
	s.Equal([]byte("test"), content)
}

func Test(t *testing.T) {
	suite.Run(t, new(FileReaderTestSuite))
}
