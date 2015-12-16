package postgres

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	_ "github.com/netw00rk/sqltractor/driver/sqlite3"
	"github.com/netw00rk/sqltractor/integration"
	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

const CONNECTION_URL = "sqlite3://integration_test.sqlite3"

var files map[string][]byte = map[string][]byte{
	"001_test.up.sql": []byte(`
CREATE TABLE test_table_1 (id INTEGER NOT NULL PRIMARY KEY);
CREATE TABLE test_table_2 (id INTEGER NOT NULL PRIMARY KEY);`),

	"001_test.down.sql": []byte(`
DROP TABLE test_table_1;
DROP TABLE test_table_2;`),

	"002_test.up.sql": []byte(`
INSERT INTO test_table_1 (id) VALUES (1);
INSERT INTO test_table_2 (id) VALUES (1);`),

	"002_test.down.sql": []byte(``),
}

type SqliteTestSuite struct {
	integration.DriverTestSuite
	connection *sql.DB
}

func (s *SqliteTestSuite) SetupSuite() {
	s.DriverTestSuite.ConnectionUrl = CONNECTION_URL
	file.SetDefaultReader(integration.NewMapFileReader(files))
}

func (s *SqliteTestSuite) TearDownSuite() {
	os.Remove("integration_test.sqlite3")
}

func TestSqliteTestSuite(t *testing.T) {
	suite.Run(t, new(SqliteTestSuite))
}
