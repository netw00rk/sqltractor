package postgres

import (
	"database/sql"
	"testing"

	"github.com/netw00rk/sqltractor/migrate/direction"
	"github.com/netw00rk/sqltractor/migrate/file"
	"github.com/netw00rk/sqltractor/migrate/pipe"
)

// TestMigrate runs some additional tests on Migrate().
// Basic testing is already done in migrate/migrate_test.go
func TestMigrate(t *testing.T) {
	driverUrl := "postgres://localhost/migratetest?sslmode=disable"

	// prepare clean database
	connection, err := sql.Open("postgres", driverUrl)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := connection.Exec(`
				DROP TABLE IF EXISTS yolo;
				DROP TABLE IF EXISTS ` + tableName + `;`); err != nil {
		t.Fatal(err)
	}

	d := &Driver{}
	if err := d.Initialize(driverUrl); err != nil {
		t.Fatal(err)
	}

	files := []file.File{
		{
			Path:      "/foobar",
			FileName:  "001_foobar.up.sql",
			Version:   1,
			Name:      "foobar",
			Direction: direction.Up,
			Content: []byte(`
				CREATE TABLE yolo (
					id serial not null primary key
				);
			`),
		},
		{
			Path:      "/foobar",
			FileName:  "002_foobar.down.sql",
			Version:   1,
			Name:      "foobar",
			Direction: direction.Down,
			Content: []byte(`
				DROP TABLE yolo;
			`),
		},
		{
			Path:      "/foobar",
			FileName:  "002_foobar.up.sql",
			Version:   2,
			Name:      "foobar",
			Direction: direction.Up,
			Content: []byte(`
				CREATE TABLE error (
					id THIS WILL CAUSE AN ERROR
				)
			`),
		},
	}

	p := pipe.New()
	go d.Migrate(files[0], p)
	errs := pipe.ReadErrors(p)
	if len(errs) > 0 {
		t.Fatal(errs)
	}

	p = pipe.New()
	go d.Migrate(files[1], p)
	errs = pipe.ReadErrors(p)
	if len(errs) > 0 {
		t.Fatal(errs)
	}

	p = pipe.New()
	go d.Migrate(files[2], p)
	errs = pipe.ReadErrors(p)
	if len(errs) == 0 {
		t.Error("Expected test case to fail")
	}

	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
}
