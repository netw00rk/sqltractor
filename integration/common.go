package integration

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/netw00rk/sqltractor/driver"
	"github.com/netw00rk/sqltractor/reader"
	"github.com/netw00rk/sqltractor/tractor"
)

type DriverTestSuite struct {
	suite.Suite
	Driver driver.Driver
	Reader reader.Reader
}

func (s *DriverTestSuite) SetupSuite() {
	fmt.Println("Starting vagrant container...")
	if _, err := exec.Command("vagrant", "up").Output(); err != nil {
		s.Fail(err.Error())
	}

	fmt.Println("Vagrant container started")
}

func (s *DriverTestSuite) TearDownSuite() {
	fmt.Println("Stoping vagrant container...")
	if _, err := exec.Command("vagrant", "halt").Output(); err != nil {
		s.Fail(err.Error())
	}
}

func (s *DriverTestSuite) TestUpAsyncDownAsync() {
	tractor := &tractor.SqlTractor{
		Driver: s.Driver,
		Reader: s.Reader,
	}

	version, _ := tractor.Version()
	s.Equal(uint64(0), version)

	for r := range tractor.UpAsync() {
		s.Nil(r.Error)
	}

	version, _ = tractor.Version()
	s.Equal(uint64(3), version)

	for r := range tractor.DownAsync() {
		s.Nil(r.Error)
	}

	version, _ = tractor.Version()
	s.Equal(uint64(0), version)
}

func (s *DriverTestSuite) TestUpDown() {
	t := &tractor.SqlTractor{
		Driver: s.Driver,
		Reader: s.Reader,
	}

	files, err := tractor.Up(t)
	s.Equal(3, len(files))
	s.Nil(err)

	version, _ := t.Version()
	s.Equal(uint64(3), version)

	files, err = tractor.Down(t)
	s.Equal(3, len(files))
	s.Nil(err)

	version, _ = t.Version()
	s.Equal(uint64(0), version)
}

func (s *DriverTestSuite) TestMigrateAsyncUpDown() {
	tractor := &tractor.SqlTractor{
		Driver: s.Driver,
		Reader: s.Reader,
	}

	version, _ := tractor.Version()
	s.Equal(uint64(0), version)

	for r := range tractor.MigrateAsync(+1) {
		s.Nil(r.Error)
	}

	version, _ = tractor.Version()
	s.Equal(uint64(1), version)

	for r := range tractor.MigrateAsync(-1) {
		s.Nil(r.Error)
	}

	version, _ = tractor.Version()
	s.Equal(uint64(0), version)
}

func (s *DriverTestSuite) TestMigrateUpDown() {
	t := &tractor.SqlTractor{
		Driver: s.Driver,
		Reader: s.Reader,
	}

	files, err := tractor.Migrate(t, +1)
	s.Equal(1, len(files))
	s.Nil(err)

	version, _ := t.Version()
	s.Equal(uint64(1), version)

	files, err = tractor.Migrate(t, -1)
	s.Equal(1, len(files))
	s.Nil(err)

	version, _ = t.Version()
	s.Equal(uint64(0), version)
}

func RepeatWhileError(fn func() error) error {
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	var err error
	for c := range ticker.C {
		if startTime.Add(60 * time.Second).Before(c) {
			return errors.New(fmt.Sprintf("Can't complite function without error %s after 60 seconds", err))
		}

		if err = fn(); err == nil {
			ticker.Stop()
			break
		}
	}

	return nil
}
