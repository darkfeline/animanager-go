package query

import "database/sql"

// Error is a superset of the error interface.
type Error interface {
	error
	// Missing returns true if the error is caused by a missing
	// row.
	Missing() bool
}

type qErr struct {
	error
}

func (e qErr) Missing() bool {
	return e.error == sql.ErrNoRows
}
