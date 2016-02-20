package tractor

import (
	"github.com/netw00rk/sqltractor/driver"
	"github.com/netw00rk/sqltractor/reader"
	"github.com/netw00rk/sqltractor/tractor/migration"
	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

// Structure for holding migration result
type Result struct {
	// Executed file
	File *file.File

	// Error
	Error error
}

// Tractor interface
type Tractor interface {
	UpAsync() chan Result
	DownAsync() chan Result
	MigrateAsync(int) chan Result
	Version() (uint64, error)
}

// SqlTractor is main structure to work with migration.
type SqlTractor struct {
	Driver driver.Driver
	Reader reader.Reader

	_manager migration.Manager
}

// Returns new pointer of SqlTractor struct or error
func NewSqlTractor(driver driver.Driver, reader reader.Reader) Tractor {
	return &SqlTractor{
		Driver: driver,
		Reader: reader,
	}
}

// Applies all available migrations asynchronously
func (t *SqlTractor) UpAsync() chan Result {
	version, err := t.Version()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	manager, err := t.manager()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	return t.applyAsync(manager.ToLastFrom(version))
}

// Rolls back all migrations asynchronously
func (t *SqlTractor) DownAsync() chan Result {
	version, err := t.Version()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	manager, err := t.manager()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	return t.applyAsync(manager.ToFirstFrom(version))
}

// Applies relative +n/-n migrations asynchronously
func (t *SqlTractor) MigrateAsync(relativeN int) chan Result {
	version, err := t.Version()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	manager, err := t.manager()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	return t.applyAsync(manager.From(version, relativeN))
}

// Returns the current migration version
func (t *SqlTractor) Version() (uint64, error) {
	driver, err := t.driver()
	if err != nil {
		return 0, err
	}

	return driver.Version()
}

func (t *SqlTractor) wrapAsyncError(err error) chan Result {
	result := make(chan Result)
	go func(result chan Result) {
		result <- Result{nil, err}
		close(result)
	}(result)
	return result
}

func (t *SqlTractor) applyAsync(files []*file.File) chan Result {
	resultChan := make(chan Result)
	go t.apply(files, resultChan)
	return resultChan
}

func (t *SqlTractor) apply(files []*file.File, resultChan chan Result) {
	if err := t.lock(); err != nil {
		resultChan <- Result{nil, err}
		close(resultChan)
		return
	}

	driver, err := t.driver()
	if err != nil {
		resultChan <- Result{nil, err}
		t.release()
		close(resultChan)
		return
	}

	for _, f := range files {
		err := driver.Migrate(f)
		resultChan <- Result{nil, err}

		if err != nil {
			t.release()
			close(resultChan)
			return
		}
	}

	t.release()
	close(resultChan)
}

func (t *SqlTractor) lock() error {
	driver, err := t.driver()
	if err != nil {
		return err
	}

	return driver.Lock()
}

func (t *SqlTractor) release() error {
	driver, err := t.driver()
	if err != nil {
		return err
	}

	return driver.Release()
}

func (t *SqlTractor) manager() (migration.Manager, error) {
	var err error
	if t._manager == nil {
		t._manager, err = migration.NewManager(t.Reader)
	}

	return t._manager, err
}

func (t *SqlTractor) driver() (driver.Driver, error) {
	if err := t.Driver.Initialize(); err != nil {
		return nil, err
	}

	return t.Driver, nil
}
