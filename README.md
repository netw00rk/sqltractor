# sqltractor

SQL schema migration tool for GO. Use it in your existing Golang code or run commands via the CLI.

## Usage from Terminal

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


## Usage in Go

See GoDoc here: http://godoc.org/github.com/mattes/migrate/migrate

```go
import "github.com/netw00rk/sqltractor/tractor"

// Import any required drivers so that they are registered and available
import _ "github.com/netw00rk/sqltractor/driver/postgres"

// use synchronous versions of migration functions ...
// files - slice of applied migration files 
files, err := tractor.Up("driver://url", "./path")
if err != nil {
  // do something with error
}

// to use the asynchronous version you have to instantiate Tractor struct
tractor, err := tractor.NewTractor("driver://url", "./path")
if err != nil {
  // do something with error
}

// UpAsync returning chan of Result struct
// type Resutl {
    File  // applied file
    Error // error if something happened
}
for r := range tractor.UpAsync() {
    if r.Error != nil {
        // do something with error
    }
    // do something with applied file
}
```
