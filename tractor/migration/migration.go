package migration

import "github.com/netw00rk/sqltractor/tractor/migration/file"

type Migration struct {
	// version of the migration file, parsed from the filenames
	Version uint64

	// reference to the *up* migration file
	UpFile *file.File

	// reference to the *down* migration file
	DownFile *file.File
}
