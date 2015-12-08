package tractor

import (
	neturl "net/url"

	"github.com/netw00rk/sqltractor/driver"
	"github.com/netw00rk/sqltractor/driver/registry"
	"github.com/netw00rk/sqltractor/tractor/file"
)

// Structure for holding migration result
type Result struct {
	// Executed file
	File *file.File

	// Error
	Error error
}

// Tractor is main structure to work with migration.
type Tractor struct {
	driver  driver.Driver
	manager file.MigrationManager
}

// Returns new pointer of Tractor struct or error
func NewTractor(url, path string) (*Tractor, error) {
	var err error

	t := new(Tractor)
	t.driver, err = initDriver(url)
	if err != nil {
		return nil, err
	}

	t.manager, err = file.NewMigrationManager(path, t.driver.FilenameExtension())
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Applies all available migrations asynchronously
func (t *Tractor) UpAsync() chan Result {
	version, err := t.driver.Version()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	return t.applyAsync(t.manager.ToLastFrom(version))
}

// Rolls back all migrations asynchronously
func (t *Tractor) DownAsync() chan Result {
	version, err := t.driver.Version()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	return t.applyAsync(t.manager.ToFirstFrom(version))
}

// Applies relative +n/-n migrations asynchronously
func (t *Tractor) MigrateAsync(relativeN int) chan Result {
	version, err := t.driver.Version()
	if err != nil {
		return t.wrapAsyncError(err)
	}

	return t.applyAsync(t.manager.From(version, relativeN))
}

// Returns the current migration version
func (t *Tractor) Version() (uint64, error) {
	return t.driver.Version()
}

func (t *Tractor) wrapAsyncError(err error) chan Result {
	result := make(chan Result)
	go func(result chan Result) {
		result <- Result{nil, err}
		close(result)
	}(result)
	return result
}

func (t *Tractor) applyAsync(files []*file.File) chan Result {
	resultChan := make(chan Result)
	go t.apply(files, resultChan)
	return resultChan
}

func (t *Tractor) apply(files []*file.File, resultChan chan Result) {
	defer t.release()
	if err := t.lock(); err != nil {
		resultChan <- Result{nil, err}
		close(resultChan)
		return
	}

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

func (t *Tractor) lock() error {
	return t.driver.Lock()
}

func (t *Tractor) release() error {
	return t.driver.Release()
}

func initDriver(url string) (driver.Driver, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}

	d, err := registry.GetDriver(u.Scheme)
	if err != nil {
		return nil, err
	}

	if err := d.Initialize(url); err != nil {
		return nil, err
	}

	return d, nil
}
