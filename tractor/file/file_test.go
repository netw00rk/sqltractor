package file

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/netw00rk/sqltractor/tractor/direction"
	"github.com/netw00rk/sqltractor/tractor/file/reader"
)

type ParserTestSuite struct {
	suite.Suite
}

type MockFileReader struct{}

func (m MockFileReader) ReadFileContent(path string) ([]byte, error) {
	return []byte("test"), nil
}

func (m MockFileReader) ReadPath(path string) ([]os.FileInfo, error) {
	return nil, nil
}

func (s *ParserTestSuite) Test() {
	var tests = []struct {
		filename          string
		expectedVersion   uint64
		expectedName      string
		expectedDirection direction.Direction
	}{
		{"001_test_file.up.sql", 1, "test_file", direction.Up},
		{"001_test_file.down.sql", 1, "test_file", direction.Down},
		{"10034_test_file.down.sql", 10034, "test_file", direction.Down},
	}

	filenameRegex := regexp.MustCompile(fmt.Sprintf(filenameRegex, "sql"))
	for _, test := range tests {
		version, name, direction, _ := parseFilenameSchema(test.filename, filenameRegex)
		s.Equal(test.expectedVersion, version, "version numbers are not equal")
		s.Equal(test.expectedName, name, "name are not equal")
		s.Equal(test.expectedDirection, direction, "directions are not equal")
	}
}

func (s *ParserTestSuite) TestReadContent() {
	SetFileReader(MockFileReader{})

	file := new(File)
	file.ReadContent()
	s.Equal([]byte("test"), file.Content)

	SetFileReader(reader.IOFileReader{})
}

func (s *ParserTestSuite) TestInvalidNames() {
	tests := []string{
		"-1_test_file.down.sql", "test_file.down.sql", "100_test_file.down",
		"100_test_file.sql", "100_test_file", "test_file", "100", ".sql", "up.sql", "down.sql"}

	filenameRegex := regexp.MustCompile(fmt.Sprintf(filenameRegex, "sql"))
	for _, test := range tests {
		_, _, _, err := parseFilenameSchema(test, filenameRegex)
		s.NotNil(err, "parsing error is nul")
	}
}

func TestParesSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
