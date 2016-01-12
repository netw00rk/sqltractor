package memory

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

var files map[string][]byte = map[string][]byte{
	"001_migrationfile.up.sql":   nil,
	"001_migrationfile.down.sql": nil,
	"002_migrationfile.up.sql":   []byte("test"),
}

type MemoryReaderTestSuite struct {
	suite.Suite
}

func (s *MemoryReaderTestSuite) TestReadFiles() {
	reader := NewMemoryReader(files)
	files, err := reader.Read()
	s.Equal(3, len(files))
	fmt.Println(files)

	var content []byte
	for _, file := range files {
		if file.FileName == "002_migrationfile.up.sql" {
			content, err = file.Content()
		}
	}

	s.Nil(err, fmt.Sprintf("Can't read %s file content", files[2].FileName))
	s.Equal([]byte("test"), content, files[2].FileName)
}

func Test(t *testing.T) {
	suite.Run(t, new(MemoryReaderTestSuite))
}
