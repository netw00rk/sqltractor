// Package cassandra implements the Driver interface.
package cassandra

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocql/gocql"

	"github.com/netw00rk/sqltractor/tractor/migration/direction"
	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

type Driver struct {
	session *gocql.Session
	url     string
}

const (
	TABLE_NAME  = "schema_migrations"
	LOCK_TABLE  = "schema_migrations_lock"
	VERSION_ROW = 1
)

// Cassandra Driver URL format:
// cassandra://host:port/keyspace
//
// Example:
// cassandra://localhost/SpaceOfKeys
func New(url string) *Driver {
	return &Driver{
		url: url,
	}
}

func FromSession(session *gocql.Session) *Driver {
	return &Driver{
		session: session,
	}
}

func (driver *Driver) Initialize() error {
	if driver.session != nil {
		return nil
	}

	u, err := url.Parse(driver.url)

	cluster := gocql.NewCluster(u.Host)
	cluster.Keyspace = u.Path[1:len(u.Path)]
	cluster.Consistency = gocql.All
	cluster.Timeout = 1 * time.Minute

	// Check if url user struct is null
	if u.User != nil {
		password, passwordSet := u.User.Password()
		if passwordSet == false {
			return fmt.Errorf("Missing password. Please provide password.")
		}

		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: u.User.Username(),
			Password: password,
		}

	}

	driver.session, err = cluster.CreateSession()
	if err != nil {
		return err
	}

	if err := driver.ensureVersionTableExists(); err != nil {
		return err
	}

	return nil
}

func (driver *Driver) Close() error {
	driver.session.Close()
	return nil
}

func (driver *Driver) Lock() error {
	if err := driver.session.Query(fmt.Sprintf("CREATE TABLE %s (lock BOOLEAN PRIMARY KEY)", LOCK_TABLE)).Exec(); err != nil {
		return err
	}

	return nil
}

func (driver *Driver) Release() error {
	if err := driver.session.Query(fmt.Sprintf("DROP TABLE %s", LOCK_TABLE)).Exec(); err != nil {
		return err
	}

	return nil
}

func (driver *Driver) Migrate(f *file.File) error {
	content, err := f.Content()
	if err != nil {
		return err
	}

	for _, query := range strings.Split(string(content), ";") {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		if err := driver.session.Query(query).Exec(); err != nil {
			return err
		}
	}

	if err := driver.version(f.Direction); err != nil {
		return err
	}

	return nil
}

func (driver *Driver) Version() (uint64, error) {
	var version int64
	err := driver.session.Query(fmt.Sprintf("SELECT version FROM %s WHERE versionRow = ?", TABLE_NAME), VERSION_ROW).Scan(&version)
	return uint64(version) - 1, err
}

func (driver *Driver) ensureVersionTableExists() error {
	err := driver.session.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (version COUNTER, versionRow BIGINT PRIMARY KEY)", TABLE_NAME)).Exec()
	if err != nil {
		return err
	}

	_, err = driver.Version()
	if err != nil {
		driver.session.Query(UP.String(), VERSION_ROW).Exec()
	}

	return nil
}

func (driver *Driver) version(d direction.Direction) error {
	var stmt counterStmt
	switch d {
	case direction.Up:
		stmt = UP
	case direction.Down:
		stmt = DOWN
	}
	return driver.session.Query(stmt.String(), VERSION_ROW).Exec()
}
