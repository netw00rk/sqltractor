// Package cassandra implements the Driver interface.
package cassandra

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocql/gocql"

	"github.com/netw00rk/sqltractor/driver/registry"
	"github.com/netw00rk/sqltractor/tractor/direction"
	"github.com/netw00rk/sqltractor/tractor/file"
)

type Driver struct {
	session *gocql.Session
}

const (
	TABLE_NAME  = "schema_migrations"
	VERSION_ROW = 1
)

// Cassandra Driver URL format:
// cassandra://host:port/keyspace
//
// Example:
// cassandra://localhost/SpaceOfKeys
func (driver *Driver) Initialize(rawurl string) error {
	u, err := url.Parse(rawurl)

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

func (driver *Driver) FilenameExtension() string {
	return "cql"
}

func (driver *Driver) Migrate(f *file.File) error {
	if err := f.ReadContent(); err != nil {
		return err
	}

	for _, query := range strings.Split(string(f.Content), ";") {
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
	err := driver.session.Query("SELECT version FROM " + TABLE_NAME + " WHERE versionRow = ?", VERSION_ROW).Scan(&version)
	return uint64(version) - 1, err
}

func (driver *Driver) ensureVersionTableExists() error {
	err := driver.session.Query("CREATE TABLE IF NOT EXISTS " + TABLE_NAME + " (version counter, versionRow bigint primary key);").Exec()
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

func init() {
	registry.RegisterDriver("cassandra", new(Driver))
}
