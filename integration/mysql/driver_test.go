package postgres

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"

	"github.com/netw00rk/sqltractor/driver/mysql"
	"github.com/netw00rk/sqltractor/integration"
	"github.com/netw00rk/sqltractor/reader/memory"
)

const CONNECTION_URL = "root@tcp(localhost:3308)/integration_test"

var files map[string][]byte = map[string][]byte{
	"001_test.up.sql": []byte(`
CREATE TABLE test_table_1 (id INT(11) NOT NULL PRIMARY KEY);
CREATE TABLE test_table_2 (id INT(11) NOT NULL PRIMARY KEY);`),

	"001_test.down.sql": []byte(`
DROP TABLE test_table_1;
DROP TABLE test_table_2;`),

	"002_test.up.sql": []byte(`
INSERT INTO test_table_1 (id) VALUES (1);
INSERT INTO test_table_2 (id) VALUES (1);`),

	"002_test.down.sql": []byte(""),
	"003_test.up.sql":   []byte(""),
	"003_test.down.sql": []byte(""),
}

type MysqlTestSuite struct {
	integration.DriverTestSuite
	connection *sql.DB
}

func (s *MysqlTestSuite) SetupSuite() {
	s.DriverTestSuite.SetupSuite()

	err := integration.RepeatWhileError(func() error {
		var err error
		if s.connection, err = sql.Open("mysql", CONNECTION_URL); err != nil {
			return err
		}
		return s.connection.Ping()
	})

	if err != nil {
		s.Fail(err.Error())
	}
}

func (s *MysqlTestSuite) SetupTest() {
	s.connection.Exec("DROP SCHEMA public CASCADE")
	s.connection.Exec("CREATE SCHEMA public")

	s.DriverTestSuite.Driver = mysql.New("mysql://" + CONNECTION_URL)
	s.DriverTestSuite.Reader = memory.NewMemoryReader(files)
}

func (s *MysqlTestSuite) TearDownSuite() {
	s.DriverTestSuite.TearDownSuite()
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}
