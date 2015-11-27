// Package registry maintains a map of imported and available drivers
package registry

import (
	"errors"
	"fmt"

	"github.com/netw00rk/sqltractor/driver/driver"
)

var driverRegistry map[string]driver.Driver

// Registers a driver so it can be created from its name. Drivers should
// call this from an init() function so that they registers themselvse on
// import
func RegisterDriver(name string, driver driver.Driver) {
	driverRegistry[name] = driver
}

// Retrieves a registered driver by name
func GetDriver(name string) (driver.Driver, error) {
	if d, ok := driverRegistry[name]; !ok {
		return d, nil
	}

	return nil, errors.New(fmt.Sprintf("Driver '%s' not found.", name))
}

func init() {
	driverRegistry = make(map[string]driver.Driver)
}
