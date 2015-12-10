package migration

import (
	"errors"
	"sort"

	"github.com/netw00rk/sqltractor/tractor/migration/direction"
	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

// Manager is a slice of Migration
type Manager []*Migration

// Initialize slice of Migration structures reads all migration files from a given path
func NewManager(path string) (Manager, error) {
	files, err := file.DefaultReader.ReadPath(path)
	if err != nil {
		return nil, err
	}

	tmp := make(map[uint64]*Migration)
	for _, file := range files {
		var migration *Migration
		var ok bool
		if migration, ok = tmp[file.Version]; !ok {
			migration = &(Migration{
				Version: file.Version,
			})
			tmp[file.Version] = migration
		}

		switch file.Direction {
		case direction.Up:
			migration.UpFile = file
		case direction.Down:
			migration.DownFile = file
		default:
			return nil, errors.New("Unsupported direction.Direction Type")
		}
	}

	newFiles := make(Manager, len(tmp))
	index := 0
	for _, v := range tmp {
		newFiles[index] = v
		index++
	}

	sort.Sort(newFiles)
	return newFiles, nil
}

// ToFirstFrom fetches all (down) migration files including the migration file
// of the current version to the very first migration file.
func (mm Manager) ToFirstFrom(version uint64) []*file.File {
	sort.Sort(sort.Reverse(mm))
	files := make([]*file.File, 0)
	for _, migration := range mm {
		if migration.Version <= version && migration.DownFile != nil {
			files = append(files, migration.DownFile)
		}
	}
	return files
}

// ToLastFrom fetches all (up) migration files to the most recent migration file.
// The migration file of the current version is not included.
func (mm Manager) ToLastFrom(version uint64) []*file.File {
	sort.Sort(mm)
	files := make([]*file.File, 0)
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
func (mm Manager) From(version uint64, relativeN int) []*file.File {
	files := make([]*file.File, 0)

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
func (mm Manager) Len() int {
	return len(mm)
}

// Less reports whether the element with
// index i should sort before the element with index j.
// Required by Sort Interface{}
func (mm Manager) Less(i, j int) bool {
	return mm[i].Version < mm[j].Version
}

// Swap swaps the elements with indexes i and j.
// Required by Sort Interface{}
func (mm Manager) Swap(i, j int) {
	mm[i], mm[j] = mm[j], mm[i]
}
