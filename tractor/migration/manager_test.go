package migration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/netw00rk/sqltractor/tractor/migration/direction"
)

type ManagerTestSuite struct {
	suite.Suite
	manager Manager
	path    string
}

func (s *ManagerTestSuite) SetupSuite() {
	s.path, _ = ioutil.TempDir("/tmp", "TestLookForMigrationFilesInSearchPath")
	ioutil.WriteFile(path.Join(s.path, "001_migrationfile.up.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "001_migrationfile.down.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "002_migrationfile.up.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "002_migrationfile.down.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "101_create_table.up.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "101_drop_tables.down.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "301_migrationfile.up.sql"), nil, 0755)
	ioutil.WriteFile(path.Join(s.path, "401_migrationfile.down.sql"), []byte("test"), 0755)

	var err error
	s.manager, err = NewManager(s.path)
	s.Nil(err, "can not read migration files")
}

func (s *ManagerTestSuite) TearDownSuite() {
	os.RemoveAll(s.path)
}

func (s *ManagerTestSuite) TestReadMigrations() {
	var tests = []struct {
		up   bool
		down bool
	}{
		{false, true}, {true, false}, {true, true}, {true, true}, {false, true},
	}

	for i, migration := range s.manager {
		if tests[i].up {
			s.NotNil(migration.UpFile, fmt.Sprintf("missing up file for version %d", migration.Version))
		}

		if tests[i].down {
			s.NotNil(migration.DownFile, fmt.Sprintf("missing down file for version %d", migration.Version))
		}
	}
}

func (s *ManagerTestSuite) TestToFirstFrom() {
	files := s.manager.ToFirstFrom(401)

	s.Equal(4, len(files), "number of files should be 4")
	for _, file := range files {
		s.True(direction.Down == file.Direction, "direction of migration should be down")
	}
}

func (s *ManagerTestSuite) TestToLastFrom() {
	files := s.manager.ToLastFrom(1)

	s.Equal(3, len(files), "nubmer of files shoud be 3")
	for _, file := range files {
		s.Equal(direction.Up, file.Direction, "direction of migration should be up")
	}
}

func (s *ManagerTestSuite) TestFrom() {
	var tests = []struct {
		from              uint64
		relative          int
		expectedVersions  []uint64
		expectedDirection direction.Direction
	}{
		{0, 2, []uint64{1, 2}, direction.Up},
		{1, 4, []uint64{2, 101, 301}, direction.Up},
		{1, 0, nil, 0},
		{0, 1, []uint64{1}, direction.Up},
		{0, 0, nil, 0},
		{101, -2, []uint64{101, 2}, direction.Down},
		{401, -1, []uint64{401}, direction.Down},
	}

	for _, test := range tests {
		files := s.manager.From(test.from, test.relative)
		s.Equal(len(files), len(test.expectedVersions))

		for i, version := range test.expectedVersions {
			s.Equal(version, files[i].Version, "migration version should be equal")
			s.Equal(test.expectedDirection, files[i].Direction, "direction of migration should be %s", test.expectedDirection)
		}
	}

}

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
