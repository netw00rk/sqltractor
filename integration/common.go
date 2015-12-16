package integration

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/stretchr/testify/suite"

	_ "github.com/netw00rk/sqltractor/driver/postgres"
	"github.com/netw00rk/sqltractor/tractor"
)

type DriverTestSuite struct {
	suite.Suite
	ConnectionUrl string
}

func (s *DriverTestSuite) SetupSuite() {
	fmt.Println("Starting vagrant container...")
	if _, err := exec.Command("vagrant", "up").Output(); err != nil {
		s.Fail(err.Error())
	}

	// waiting for database to start
	time.Sleep(5 * time.Second)
	fmt.Println("Vagrant container started")
}

func (s *DriverTestSuite) TearDownSuite() {
	fmt.Println("Stoping vagrant container...")
	if _, err := exec.Command("vagrant", "halt").Output(); err != nil {
		s.Fail(err.Error())
	}
}

func (s *DriverTestSuite) TestUpAsyncDownAsync() {
	tractor, _ := tractor.NewTractor(s.ConnectionUrl, "/some/path/to/files")
	version, _ := tractor.Version()
	s.Equal(uint64(0), version)

	for r := range tractor.UpAsync() {
		s.Nil(r.Error)
	}

	version, _ = tractor.Version()
	s.Equal(uint64(2), version)

	for r := range tractor.DownAsync() {
		s.Nil(r.Error)
	}

	version, _ = tractor.Version()
	s.Equal(uint64(0), version)
}

func (s *DriverTestSuite) TestUpDown() {
	files, err := tractor.Up(s.ConnectionUrl, "/some/path/to/files")
	s.Equal(2, len(files))
	s.Nil(err)

	version, _ := tractor.Version(s.ConnectionUrl, "/some/path/to/files")
	s.Equal(uint64(2), version)

	files, err = tractor.Down(s.ConnectionUrl, "/some/path/to/files")
	s.Equal(2, len(files))
	s.Nil(err)

	version, _ = tractor.Version(s.ConnectionUrl, "/some/path/to/files")
	s.Equal(uint64(0), version)
}

func (s *DriverTestSuite) TestMigrateAsyncUpDown() {
	tractor, err := tractor.NewTractor(s.ConnectionUrl, "/some/path/to/files")
	if err != nil {
		s.Fail(err.Error())
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
	files, err := tractor.Migrate(s.ConnectionUrl, "/some/path/to/files", +1)
	s.Equal(1, len(files))
	s.Nil(err)

	version, _ := tractor.Version(s.ConnectionUrl, "/some/path/to/files")
	s.Equal(uint64(1), version)

	files, err = tractor.Migrate(s.ConnectionUrl, "/some/path/to/files", -1)
	s.Equal(1, len(files))
	s.Nil(err)

	version, _ = tractor.Version(s.ConnectionUrl, "/some/path/to/files")
	s.Equal(uint64(0), version)
}
