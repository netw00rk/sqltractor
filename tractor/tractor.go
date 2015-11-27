package tractor

import (
	utils "github.com/netw00rk/sqltractor/driver"
	"github.com/netw00rk/sqltractor/driver/driver"
	"github.com/netw00rk/sqltractor/tractor/file"
)

// Struct for holding migration resutl
type Result struct {
	// Executed file
	File *file.File

	// Error
	Error error
}

// Tractor is main struct to work with migration
// To instantiate use NewTractor func
type Tractor struct {
	driver driver.Driver
	manager file.MigrationManager
}

// Returns new pointer of Tractor struct or error if
// url format is driver://user:password@host/database
// path is filesystem path with migration files
func NewTractor(url, path string) (*Tractor, error) {
	var err error

	t := new(Tractor)
	t.driver, err = utils.New(url)
	if err != nil {
		return nil, err
	}

	t.manager, err = file.NewMigrationManager(path, t.driver.FilenameExtension())
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Up applies all available migrations
func (t *Tractor) UpAsync() chan Result {
	version, err := t.driver.Version()
	if err != nil {
		panic(err)
	}

	return t.applyAsync(t.manager.ToLastFrom(version))
}

// Down rolls back all migrations
func (t *Tractor) DownAsync() chan Result {
	version, err := t.driver.Version()
	if err != nil {
		panic(err)
	}

	return t.applyAsync(t.manager.ToFirstFrom(version))
}

// Migrate applies relative +n/-n migrations
func (t *Tractor) MigrateAsync(relativeN int) chan Result {
	version, err := t.driver.Version()
	if err != nil {
		panic(err)
	}

	return t.applyAsync(t.manager.From(version, relativeN))
}

// Version returns the current migration version
func (t *Tractor) Version() (uint64, error) {
	return t.driver.Version()
}

func (t *Tractor) applyAsync(files []*file.File) chan Result {
	resultChan := make(chan Result)
	go t.apply(files, resultChan)
	return resultChan
}

func (t *Tractor) apply(files []*file.File, resultChan chan Result) {
	for _, f := range files {
		err := t.driver.Migrate(f)
		resultChan <- Result{nil, err}

		if err != nil {
			close(resultChan)
			return
		}
	}
	close(resultChan)
}
