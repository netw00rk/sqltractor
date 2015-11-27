package tractor

import "github.com/netw00rk/sqltractor/tractor/file"

func Up(url, path string) ([]*file.File, error) {
	files := make([]*file.File, 0)

	t, err := NewTractor(url, path)
	if err != nil {
		return files, err
	}

	for r := range t.UpAsync() {
		if r.Error != nil {
			return files, err
		}
		files = append(files, r.File)
	}

	return files, nil
}

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

func Version(url, path string) (version uint64, err error) {
	t, err := NewTractor(url, path)
	if err != nil {
		return 0, err
	}
	return t.Version()
}