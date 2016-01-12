# sqltractor

SQL schema migration tool for Go. Use it in your existing Go code or run commands via the CLI.

## Usage

**from Terminal**

```bash
# install
go get github.com/netw00rk/sqltractor/sqltractor-cli

# create new migration file in path
sqltractor-cli -url driver://url -path ./migrations create migration_file_xyz

# apply all available migrations
sqltractor-cli -url driver://url -path ./migrations up

# roll back all migrations
sqltractor-cli -url driver://url -path ./migrations down

# show the current migration version
sqltractor-cli -url driver://url -path ./migrations version

# apply the next n migrations
sqltractor-cli -url driver://url -path ./migrations migrate +1
sqltractor-cli -url driver://url -path ./migrations migrate +2
sqltractor-cli -url driver://url -path ./migrations migrate +n

# roll back the previous n migrations
sqltractor-cli -url driver://url -path ./migrations migrate -1
sqltractor-cli -url driver://url -path ./migrations migrate -2
sqltractor-cli -url driver://url -path ./migrations migrate -n

# go to specific migration
sqltractor-cli -url driver://url -path ./migrations goto 1
sqltractor-cli -url driver://url -path ./migrations goto 10
sqltractor-cli -url driver://url -path ./migrations goto v
```

**in Go code**

See GoDoc here: http://godoc.org/github.com/netw00rk/sqltractor/tractor

```go
import "github.com/netw00rk/sqltractor/tractor"

// Import required driver and reader
import "github.com/netw00rk/sqltractor/driver/postgres"
import "github.com/netw00rk/sqltractor/reader/file"


func main() {
    // create tractor struct
    t := &tractor.SqlTractor{
        Driver: driver.New("driver://url")
        Reader: file.NewFileReader("./path/to/migration/files")
    }

    // UpAsync returning chan of Result structure
    // type Result struct {
    //    File  // applied file
    //    Error // error if something happened
    //}
    for r := range tractor.UpAsync() {
        if r.Error != nil {
            if r.File != nil {
                fmt.Printf("Error %s while applying file %s", r.Error, r.File.FileName)
            }
        }
        // do something with applied file
    }

    // usage of synchronous wrapper that
    // return slice of applied migration files and error
    files, err := tractor.Up(t)
    if err != nil {
      // do something with error
    }
}
```

## Available Drivers

 * [PostgreSQL](https://github.com/netw00rk/sqltractor/tree/master/driver/postgres)
 * [Cassandra](https://github.com/netw00rk/sqltractor/tree/master/driver/cassandra)
 * [SQLite](https://github.com/netw00rk/sqltractor/tree/master/driver/sqlite3)
 * [MySQL](https://github.com/netw00rk/sqltractor/tree/master/driver/mysql) (experimental)

Need another driver? Just implement the [Driver interface](http://godoc.org/github.com/netw00rk/sqltractor/driver#Driver) and open a PR.

## Avaiable readers

 * [FileReader](https://github.com/netw00rk/sqltractor/tree/master/reader/file)
 * [MemoryReader](https://github.com/netw00rk/sqltractor/tree/master/reader/memory)

## Migration files

The format of migration files looks like this:

```
001_initial_plan_to_do_sth.up.sql     # up migration instructions
001_initial_plan_to_do_sth.down.sql   # down migration instructions
002_xxx.up.sql
002_xxx.down.sql
...
```

Why two files? This way you could still do sth like 
``psql -f ./db/migrations/001_initial_plan_to_do_sth.up.sql`` and there is no
need for any custom markup language to divide up and down migrations. Please note
that the filename extension depends on the driver.


## Acknowledgements

Many thanks goes to Matthias Kadenbach, https://github.com/mattes and all contributors to the https://github.com/mattes/migrate for the ideas and code
