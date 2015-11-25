// Package postgres implements the Driver interface.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/lib/pq"

	"github.com/netw00rk/sqltractor/driver/registry"
	"github.com/netw00rk/sqltractor/migrate/direction"
	"github.com/netw00rk/sqltractor/migrate/file"
)

type Driver struct {
	db *sql.DB
}

const TABLE_NAME string = "schema_migrations"
var tableName string

func (driver *Driver) Initialize(rawurl string) error {
	db, err := sql.Open("postgres", rawurl)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	driver.db = db
	if err := driver.ensureSchemaExists(extractCurrentSchema(rawurl), extractUser(rawurl)); err != nil {
		return err
	}

	if err := driver.ensureVersionTableExists(); err != nil {
		return err
	}
	return nil
}

func (driver *Driver) Close() error {
	if err := driver.db.Close(); err != nil {
		return err
	}
	return nil
}

func (driver *Driver) ensureSchemaExists(schema, user string) error {
	if schema != "" {
		if _, err := driver.db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)); err != nil {
			return err
		}
		tableName = fmt.Sprintf("%s.%s", schema, TABLE_NAME)
	}
	return nil
}

func (driver *Driver) ensureVersionTableExists() error {
	if _, err := driver.db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (version int not null primary key);"); err != nil {
		return err
	}
	return nil
}

func (driver *Driver) FilenameExtension() string {
	return "sql"
}

func (driver *Driver) Migrate(f file.File, pipe chan interface{}) {
	defer close(pipe)
	pipe <- f

	tx, err := driver.db.Begin()
	if err != nil {
		pipe <- err
		return
	}

	if f.Direction == direction.Up {
		if _, err := tx.Exec("INSERT INTO "+tableName+" (version) VALUES ($1)", f.Version); err != nil {
			pipe <- err
			if err := tx.Rollback(); err != nil {
				pipe <- err
			}
			return
		}
	} else if f.Direction == direction.Down {
		if _, err := tx.Exec("DELETE FROM "+tableName+" WHERE version=$1", f.Version); err != nil {
			pipe <- err
			if err := tx.Rollback(); err != nil {
				pipe <- err
			}
			return
		}
	}

	if err := f.ReadContent(); err != nil {
		pipe <- err
		return
	}

	if _, err := tx.Exec(string(f.Content)); err != nil {
		fmt.Println(err)
		pqErr := err.(*pq.Error)
		offset, err := strconv.Atoi(pqErr.Position)
		if err == nil && offset >= 0 {
			lineNo, columnNo := file.LineColumnFromOffset(f.Content, offset-1)
			errorPart := file.LinesBeforeAndAfter(f.Content, lineNo, 5, 5, true)
			pipe <- errors.New(fmt.Sprintf("%s %v: %s in line %v, column %v:\n\n%s", pqErr.Severity, pqErr.Code, pqErr.Message, lineNo, columnNo, string(errorPart)))
		} else {
			pipe <- errors.New(fmt.Sprintf("%s %v: %s", pqErr.Severity, pqErr.Code, pqErr.Message))
		}

		if err := tx.Rollback(); err != nil {
			pipe <- err
		}
		return
	}

	if err := tx.Commit(); err != nil {
		pipe <- err
		return
	}
}

func (driver *Driver) Version() (uint64, error) {
	var version uint64
	err := driver.db.QueryRow("SELECT version FROM " + tableName + " ORDER BY version DESC LIMIT 1").Scan(&version)
	switch {
	case err == sql.ErrNoRows:
		return 0, nil
	case err != nil:
		return 0, err
	default:
		return version, nil
	}
}

func extractCurrentSchema(rawurl string) string {
	u, _ := url.Parse(rawurl);
	search_path := u.Query().Get("search_path")
	parts := strings.Split(search_path, ",")
	if len(parts) > 0 {
		return strings.Trim(parts[0], " ")
	}
	return ""
}

func extractUser(rawurl string) string {
	u, _ := url.Parse(rawurl);
	return u.User.Username()
}

func init() {
	registry.RegisterDriver("postgres", Driver{})
}
