package postgres

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/suite"

	"github.com/netw00rk/sqltractor/driver/cassandra"
	"github.com/netw00rk/sqltractor/integration"
	"github.com/netw00rk/sqltractor/reader/memory"
)

const CONNECTION_URL = "cassandra://localhost/integration_test"

var files map[string][]byte = map[string][]byte{
	"001_test.up.sql": []byte(`
CREATE TABLE test_table_1 (id VARINT PRIMARY KEY, msg TEXT);
CREATE INDEX ON test_table_1 (msg);`),

	"001_test.down.sql": []byte(`DROP TABLE test_table_1;`),

	"002_test.up.sql": []byte(`
INSERT INTO test_table_1 (id, msg) VALUES (1, 'some_text');
INSERT INTO test_table_1 (id, msg) VALUES (2, 'some_more_text');`),

	"002_test.down.sql": []byte(""),
	"003_test.up.sql":   []byte(""),
	"003_test.down.sql": []byte(""),
}

type CassandraTestSuite struct {
	integration.DriverTestSuite
	session *gocql.Session
}

func (s *CassandraTestSuite) SetupSuite() {
	s.DriverTestSuite.SetupSuite()
	cluster := gocql.NewCluster("localhost")
	cluster.Consistency = gocql.All
	cluster.Timeout = 1 * time.Minute

	err := integration.RepeatWhileError(func() error {
		var err error
		s.session, err = cluster.CreateSession()
		return err
	})
	s.Nil(err)

	s.session.Query("CREATE KEYSPACE integration_test WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }").Exec()
}

func (s *CassandraTestSuite) SetupTest() {
	s.DriverTestSuite.Driver = cassandra.New(CONNECTION_URL)
	s.DriverTestSuite.Reader = memory.NewMemoryReader(files)
}

func (s *CassandraTestSuite) TearDownSuite() {
	s.session.Query("DROP KEYSPACE integration_test").Exec()
	s.DriverTestSuite.TearDownSuite()
}

func TestCassandraTestSuite(t *testing.T) {
	suite.Run(t, new(CassandraTestSuite))
}

func init() {
	log.SetOutput(ioutil.Discard)
}
