package tractor

import (
	"github.com/netw00rk/sqltractor/tractor/migration/file"
)

// Utility function that wraps Tractor structure init and call UpAsync().
// Returns slice of applied files and error if something happened.
func Up(url, path string) ([]*file.File, error) {
	files := make([]*file.File, 0)

	t, err := NewTractor(url, path)
	if err != nil {
		return files, err
	}

	for r := range t.UpAsync() {
		if r.Error != nil {
			return files, r.Error
		}
		files = append(files, r.File)
	}

	return files, nil
}

// Utility function that wraps Tractor structure init and call DownAsync().
// Returns slice of applied files and error if something happened.
func Down(url, path string) ([]*file.File, error) {
	files := make([]*file.File, 0)

	t, err := NewTractor(url, path)
	if err != nil {
		return files, err
	}

	for r := range t.DownAsync() {
		if r.Error != nil {
			return files, err
		}
		files = append(files, r.File)
	}

	return files, nil
}

// Utility function that wraps Tractor structure init and call MigrateAsync().
// Returns slice of applied files and error if something happened.
func Migrate(url, path string, relativeN int) ([]*file.File, error) {
	files := make([]*file.File, 0)

	t, err := NewTractor(url, path)
	if err != nil {
		return files, err
	}

	for r := range t.MigrateAsync(relativeN) {
		if r.Error != nil {
			return files, err
		}
		files = append(files, r.File)
	}

	return files, nil
}

// Utility that wraps Tractor structure init and call Version().
// Returns current applied version and error if something happened.
func Version(url, path string) (version uint64, err error) {
	t, err := NewTractor(url, path)
	if err != nil {
		return 0, err
	}
	return t.Version()
}
