package reader

import "github.com/netw00rk/sqltractor/tractor/migration/file"

type Reader interface {
	// function that reader migration files, returns slice of File struct or error
	Read() ([]*file.File, error)
}
