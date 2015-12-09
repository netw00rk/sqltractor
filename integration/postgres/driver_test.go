package postgres

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/netw00rk/sqltractor/tractor"
	_ "github.com/netw00rk/sqltractor/driver/postgres"
)

type PostgresTestSuite struct {
	suite.Suite
}

type MockFileReader struct {
	files map[string] {
		"001_test.up.sql": `CREATE TABLE test_table (id INTEGER PRIMARY KEY);
                            CREATE TABLE test_table_2 (id INTEGER PRYMARY KEY);`,
	}
}

func (m MockFileReader) ReadFileContent(path string) ([]byte, error) {
	return []byte("test"), nil
}

func (m MockFileReader) ReadPath(path string) ([]os.FileInfo, error) {
	return nil, nil
}
