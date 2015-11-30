package file

import (
	"sort"
	"errors"
	"regexp"
	"fmt"

	"github.com/netw00rk/sqltractor/tractor/direction"
	"github.com/netw00rk/sqltractor/tractor/file/reader"
)

var (
	filenameRegex = `^([0-9]+)_(.*)\.(up|down)\.%s$`
	defaultFileReader reader.FileReader
)

func SetFileReader(r reader.FileReader) {
	defaultFileReader = r
}

type Migration struct {
	// version of the migration file, parsed from the filenames
	Version uint64

	// reference to the *up* migration file
	UpFile *File

	// reference to the *down* migration file
	DownFile *File
}

// MigrationManager is a slice of Migration
type MigrationManager []*Migration

// Initialize slice of Migration structures reads all migration files from a given path
func NewMigrationManager(path, extension string) (MigrationManager, error) {
	ioFiles, err := defaultFileReader.ReadPath(path)
	if err != nil {
		return nil, err
	}

	filenameRegex := regexp.MustCompile(fmt.Sprintf(filenameRegex, extension))
	tmp := make(map[uint64]*Migration)
	for _, file := range ioFiles {
		version, name, d, err := parseFilenameSchema(file.Name(), filenameRegex)
		if err == nil {
			var migrationFile *Migration
			var ok bool
			if migrationFile, ok = tmp[version]; !ok {
				migrationFile = &(Migration{
					Version: version,
				})
				tmp[version] = migrationFile
			}

			switch d {
			case direction.Up:
				migrationFile.UpFile = &File{
					Path:      path,
					FileName:  file.Name(),
					Version:   version,
					Name:      name,
					Content:   nil,
					Direction: direction.Up,
				}
			case direction.Down:
				migrationFile.DownFile = &File{
					Path:      path,
					FileName:  file.Name(),
					Version:   version,
					Name:      name,
					Content:   nil,
					Direction: direction.Down,
				}
			default:
				return nil, errors.New("Unsupported direction.Direction Type")
			}
		}
	}

	newFiles := make(MigrationManager, len(tmp))
	index := 0;
	for _, v := range tmp {
		newFiles[index] = v
		index++
	}

	sort.Sort(newFiles)
	return newFiles, nil
}


// ToFirstFrom fetches all (down) migration files including the migration file
// of the current version to the very first migration file.
func (mm MigrationManager) ToFirstFrom(version uint64) []*File {
	sort.Sort(sort.Reverse(mm))
	files := make([]*File, 0)
	for _, migration := range mm {
		if migration.Version <= version && migration.DownFile != nil {
			files = append(files, migration.DownFile)
		}
	}
	return files
}

// ToLastFrom fetches all (up) migration files to the most recent migration file.
// The migration file of the current version is not included.
func (mm MigrationManager) ToLastFrom(version uint64) []*File {
	sort.Sort(mm)
	files := make([]*File, 0)
	for _, migration := range mm {
		if migration.Version > version && migration.UpFile != nil {
			files = append(files, migration.UpFile)
		}
	}
	return files
}

// From travels relatively through migration files.
//
// 		+1 will fetch the next up migration file
// 		+2 will fetch the next two up migration files
// 		+n will fetch ...
// 		-1 will fetch the the previous down migration file
// 		-2 will fetch the next two previous down migration files
//		-n will fetch ...
func (mm MigrationManager) From(version uint64, relativeN int) []*File {
	files := make([]*File, 0)

	var d direction.Direction
	if relativeN > 0 {
		d = direction.Up
	} else if relativeN < 0 {
		d = direction.Down
	} else { // relativeN == 0
		return files
	}

	if d == direction.Down {
		sort.Sort(sort.Reverse(mm))
	} else {
		sort.Sort(mm)
	}

	counter := relativeN
	if relativeN < 0 {
		counter = relativeN * -1
	}

	for _, migration := range mm {
		if counter > 0 {
			if d == direction.Up && migration.Version > version && migration.UpFile != nil {
				files = append(files, migration.UpFile)
				counter -= 1
			} else if d == direction.Down && migration.Version <= version && migration.DownFile != nil {
				files = append(files, migration.DownFile)
				counter -= 1
			}
		} else {
			break
		}
	}
	return files
}

// Len is the number of elements in the collection.
// Required by Sort Interface{}
func (mm MigrationManager) Len() int {
	return len(mm)
}

// Less reports whether the element with
// index i should sort before the element with index j.
// Required by Sort Interface{}
func (mm MigrationManager) Less(i, j int) bool {
	return mm[i].Version < mm[j].Version
}

// Swap swaps the elements with indexes i and j.
// Required by Sort Interface{}
func (mm MigrationManager) Swap(i, j int) {
	mm[i], mm[j] = mm[j], mm[i]
}

func init() {
	SetFileReader(reader.IOFileReader{})
}