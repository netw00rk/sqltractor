package mysql

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/netw00rk/sqltractor/migrate/direction"
	"github.com/netw00rk/sqltractor/migrate/file"
	"github.com/netw00rk/sqltractor/migrate/pipe"
)

// TestMigrate runs some additional tests on Migrate().
// Basic testing is already done in migrate/migrate_test.go
func TestMigrate(t *testing.T) {
	driverUrl := "mysql://root@tcp(127.0.0.1:3306)/migratetest"

	// prepare clean database
	connection, err := sql.Open("mysql", strings.SplitN(driverUrl, "mysql://", 2)[1])
	if err != nil {
		t.Fatal(err)
	}

	if _, err := connection.Exec(`DROP TABLE IF EXISTS yolo, yolo1, ` + tableName); err != nil {
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
          id int(11) not null primary key auto_increment
        );

				CREATE TABLE yolo1 (
				  id int(11) not null primary key auto_increment
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

      	// a comment
				CREATE TABLE error (
          id THIS WILL CAUSE AN ERROR
        );
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
