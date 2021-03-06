// Package driver holds the driver interface.
package driver

import "github.com/netw00rk/sqltractor/tractor/migration/file"

// Driver is the interface type that needs to implemented by all drivers.
type Driver interface {

	// Initialize is the first function to be called.
	// Check the url string and open and verify any connection
	// that has to be made.
	Initialize() error

	// Close is the last function to be called.
	// Close any open connection here.
	Close() error

	// Migrate is the heart of the driver.
	// It will receive a file which the driver should apply
	// to its backend or whatever.
	Migrate(file *file.File) error

	// Version returns the current migration version.
	Version() (uint64, error)

	// Lock create lock table
	Lock() error

	// Release drops a lock table
	Release() error
}
