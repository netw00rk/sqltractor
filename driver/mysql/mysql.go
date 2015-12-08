// Package mysql implements the Driver interface.
package mysql

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/netw00rk/sqltractor/driver/registry"
	"github.com/netw00rk/sqltractor/tractor/direction"
	"github.com/netw00rk/sqltractor/tractor/file"
)

type Driver struct {
	db *sql.DB
}

const (
	TABLE_NAME      = "schema_migrations"
	LOCK_TABLE_NAME = "schema_migrations_lock"
)

var errRegexp, _ = regexp.Compile(`at line ([0-9]+)$`)

func (driver *Driver) Initialize(url string) error {
	urlWithoutScheme := strings.SplitN(url, "mysql://", 2)
	if len(urlWithoutScheme) != 2 {
		return errors.New("invalid mysql:// scheme")
	}

	db, err := sql.Open("mysql", urlWithoutScheme[1])
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
	if _, err := driver.db.Exec(fmt.Sprintf("CREATE TABLE %s (lock BOOLEAN NOT NULL)", LOCK_TABLE_NAME)); err != nil {
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
	// http://go-database-sql.org/modifying.html, Working with Transactions
	// You should not mingle the use of transaction-related functions such as Begin() and Commit() with SQL statements such as BEGIN and COMMIT in your SQL code.
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
		if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE version = ?", TABLE_NAME), f.Version); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := f.ReadContent(); err != nil {
		return err
	}

	// TODO this is not good! unfortunately there is no mysql driver that
	// supports multiple statements per query.
	sqlStmts := bytes.Split(f.Content, []byte(";"))
	for _, sqlStmt := range sqlStmts {
		sqlStmt = bytes.TrimSpace(sqlStmt)
		if len(sqlStmt) > 0 {
			if _, err := tx.Exec(string(sqlStmt)); err != nil {
				if mysqlErr, ok := err.(*mysql.MySQLError); ok {
					var lineNo int
					lineNoRe := errRegexp.FindStringSubmatch(mysqlErr.Message)
					if len(lineNoRe) == 2 {
						lineNo, err = strconv.Atoi(lineNoRe[1])
					}
					if err == nil {
						// get white-space offset
						// TODO this is broken, because we use sqlStmt instead of f.Content
						wsLineOffset := 0
						b := bufio.NewReader(bytes.NewBuffer(sqlStmt))
						for {
							line, _, err := b.ReadLine()
							if err != nil {
								break
							}
							if bytes.TrimSpace(line) == nil {
								wsLineOffset += 1
							} else {
								break
							}
						}

						message := mysqlErr.Error()
						message = errRegexp.ReplaceAllString(message, fmt.Sprintf("at line %v", lineNo+wsLineOffset))

						errorPart := file.LinesBeforeAndAfter(sqlStmt, lineNo, 5, 5, true)
						return errors.New(fmt.Sprintf("%s\n\n%s", message, string(errorPart)))
					} else {
						return errors.New(mysqlErr.Error())
					}

					if err := tx.Rollback(); err != nil {
						return err
					}

					return nil
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (driver *Driver) Version() (uint64, error) {
	var version uint64
	err := driver.db.QueryRow(fmt.Sprintf("SELECT version FROM %s ORDER BY version DESC", TABLE_NAME)).Scan(&version)
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
	_, err := driver.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (version INT NOT NULL PRIMARY KEY)", TABLE_NAME))
	if _, isWarn := err.(mysql.MySQLWarnings); err != nil && !isWarn {
		return err
	}

	return nil
}

func init() {
	registry.RegisterDriver("mysql", new(Driver))
}
