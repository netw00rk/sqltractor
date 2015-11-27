// Package driver holds the driver interface.
package driver

import (
	neturl "net/url" // alias to allow `url string` func signature in New

	"github.com/netw00rk/sqltractor/driver/registry"
	"github.com/netw00rk/sqltractor/driver/driver"
)

// New returns Driver and calls Initialize on it
func New(url string) (driver.Driver, error) {
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
