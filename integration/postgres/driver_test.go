package postgres

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/netw00rk/sqltractor/driver/postgres"
	"github.com/netw00rk/sqltractor/integration"
	"github.com/netw00rk/sqltractor/reader/memory"
)

const CONNECTION_URL = "postgres://postgres@localhost:6032/integration_test?sslmode=disable"

var maxVersion = 3
var files map[string][]byte = map[string][]byte{
	"001_test.up.sql": []byte(`
CREATE TYPE TEST_TYPE AS ENUM('value');
CREATE TABLE test_table_1 (id INTEGER NOT NULL PRIMARY KEY);
CREATE TABLE test_table_2 (id INTEGER NOT NULL PRIMARY KEY);`),

	"001_test.down.sql": []byte(`
DROP TYPE TEST_TYPE;
DROP TABLE test_table_1;
DROP TABLE test_table_2;`),

	"002_test.up.sql": []byte(`
INSERT INTO test_table_1 (id) VALUES (1);
INSERT INTO test_table_2 (id) VALUES (1);`),

	"002_test.down.sql": []byte(``),

	"003_test.up.sql": []byte(`
-- tag:no_transaction
ALTER TYPE TEST_TYPE ADD VALUE 'new_value';`),

	"003_test.down.sql": []byte(``),
}

type PostgresTestSuite struct {
	integration.DriverTestSuite
	connection *sql.DB
}

func (s *PostgresTestSuite) SetupSuite() {
	s.DriverTestSuite.SetupSuite()

	err := integration.RepeatWhileError(func() error {
		var err error
		if s.connection, err = sql.Open("postgres", CONNECTION_URL); err != nil {
			return err
		}
		return s.connection.Ping()
	})

	if err != nil {
		s.Fail(err.Error())
	}
}

func (s *PostgresTestSuite) SetupTest() {
	s.connection.Exec("DROP SCHEMA public CASCADE")
	s.connection.Exec("CREATE SCHEMA public")

	s.DriverTestSuite.Driver = postgres.New(CONNECTION_URL)
	s.DriverTestSuite.Reader = memory.NewMemoryReader(files)
}

func (s *PostgresTestSuite) TearDownSuite() {
	s.DriverTestSuite.TearDownSuite()
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
