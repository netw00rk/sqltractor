// Package sqlite3 implements the Driver interface.
package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"

	"github.com/netw00rk/sqltractor/driver/registry"
	"github.com/netw00rk/sqltractor/tractor/direction"
	"github.com/netw00rk/sqltractor/tractor/file"
)

type Driver struct {
	db *sql.DB
}

const (
	TABLE_NAME      = "schema_migration"
	LOCK_TABLE_NAME = "schema_migration_lock"
)

func (driver *Driver) Initialize(url string) error {
	filename := strings.SplitN(url, "sqlite3://", 2)
	if len(filename) != 2 {
		return errors.New("invalid sqlite3:// scheme")
	}

	db, err := sql.Open("sqlite3", filename[1])
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	driver.db = db

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

func (driver *Driver) FilenameExtension() string {
	return "sql"
}

func (driver *Driver) Lock() error {
	if _, err := driver.db.Exec(fmt.Sprintf("CREATE TABLE %s (lock INTEGER NOT NULL);", LOCK_TABLE_NAME)); err != nil {
		return err
	}

	return nil
}

func (driver *Driver) Release() error {
	if _, err := driver.db.Exec(fmt.Sprintf("DROP TABLE %s", LOCK_TABLE_NAME)); err != nil {
		return err
	}

	return nil
}

func (driver *Driver) Migrate(f *file.File) error {
	tx, err := driver.db.Begin()
	if err != nil {
		return err
	}

	if f.Direction == direction.Up {
		if _, err := tx.Exec(fmt.Sprintf("INSERT INTO %s (version) VALUES (?)", TABLE_NAME), f.Version); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	} else if f.Direction == direction.Down {
		if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE version=?", TABLE_NAME), f.Version); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := f.ReadContent(); err != nil {
		return err
	}

	if _, err := tx.Exec(string(f.Content)); err != nil {
		if sqliteErr, isErr := err.(sqlite3.Error); isErr {
			// The sqlite3 library only provides error codes, not position information. Output what we do know
			return errors.New(fmt.Sprintf("SQLite Error (%s); Extended (%s)\nError: %s", sqliteErr.Code.Error(), sqliteErr.ExtendedCode.Error(), sqliteErr.Error()))
		} else {
			return errors.New(fmt.Sprintf("An error occurred: %s", err.Error()))
		}

		if err := tx.Rollback(); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (driver *Driver) Version() (uint64, error) {
	var version uint64
	err := driver.db.QueryRow(fmt.Sprintf("SELECT version FROM %s ORDER BY version DESC LIMIT 1", TABLE_NAME)).Scan(&version)
	switch {
	case err == sql.ErrNoRows:
		return 0, nil
	case err != nil:
		return 0, err
	default:
		return version, nil
	}
}

func (driver *Driver) ensureVersionTableExists() error {
	if _, err := driver.db.Exec("CREATE TABLE IF NOT EXISTS " + TABLE_NAME + " (version INTEGER PRIMARY KEY AUTOINCREMENT);"); err != nil {
		return err
	}
	return nil
}

func init() {
	registry.RegisterDriver("sqlite3", new(Driver))
}
