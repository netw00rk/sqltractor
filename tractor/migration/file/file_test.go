package file

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/netw00rk/sqltractor/tractor/migration/direction"
)

type ParserTestSuite struct {
	suite.Suite
}

func MockedContentFunc() ([]byte, error) {
	return []byte("test"), nil
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

	for _, test := range tests {
		file, err := NewFile(test.filename, MockedContentFunc)
		s.Nil(err, "can't parse filename")
		s.Equal(test.expectedVersion, file.Version, "version numbers are not equal")
		s.Equal(test.expectedName, file.Name, "name are not equal")
		s.Equal(test.expectedDirection, file.Direction, "directions are not equal")
	}
}

func (s *ParserTestSuite) TestReadContent() {
	file := new(File)
	file.ContentFunc = MockedContentFunc
	content, _ := file.Content()
	s.Equal([]byte("test"), content)
}

func (s *ParserTestSuite) TestInvalidNames() {
	tests := []string{
		"-1_test_file.down.sql", "test_file.down.sql", "100_test_file.down",
		"100_test_file.sql", "100_test_file", "test_file", "100", ".sql", "up.sql", "down.sql"}

	for _, test := range tests {
		_, err := NewFile(test, MockedContentFunc)
		s.NotNil(err, "parsing error is nul")
	}
}

func TestParesSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
