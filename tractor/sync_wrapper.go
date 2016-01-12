package tractor

import (
	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

func Up(t Tractor) ([]*file.File, error) {
	files := make([]*file.File, 0)
	for r := range t.UpAsync() {
		if r.Error != nil {
			return files, r.Error
		}
		files = append(files, r.File)
	}

	return files, nil
}

func Down(t Tractor) ([]*file.File, error) {
	files := make([]*file.File, 0)
	for r := range t.DownAsync() {
		if r.Error != nil {
			return files, r.Error
		}
		files = append(files, r.File)
	}

	return files, nil
}

func Migrate(t Tractor, relativeN int) ([]*file.File, error) {
	files := make([]*file.File, 0)
	for r := range t.MigrateAsync(relativeN) {
		if r.Error != nil {
			return files, r.Error
		}
		files = append(files, r.File)
	}

	return files, nil
}
