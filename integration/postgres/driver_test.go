package postgres

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	_ "github.com/netw00rk/sqltractor/driver/postgres"
	"github.com/netw00rk/sqltractor/integration"
	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

const CONNECTION_URL = "postgres://postgres@localhost:6032/integration_test?sslmode=disable"

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

type PostgresTestSuite struct {
	integration.DriverTestSuite
	connection *sql.DB
}

func (s *PostgresTestSuite) SetupSuite() {
	s.DriverTestSuite.SetupSuite()

	s.connection, _ = sql.Open("postgres", CONNECTION_URL)
	s.DriverTestSuite.ConnectionUrl = CONNECTION_URL
	file.SetDefaultReader(file.NewMemoryReader(files))
}

func (s *PostgresTestSuite) SetupTest() {
	s.connection.Exec("DROP SCHEMA public CASCADE")
	s.connection.Exec("CREATE SCHEMA public")
}

func (s *PostgresTestSuite) TearDownSuite() {
	s.DriverTestSuite.TearDownSuite()
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
