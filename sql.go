package construct

import "errors"

// ErrNotFound will be returned from generated scan functions if the query result was empty, but a row was scanned
var ErrNotFound = errors.New("record not found")

// RowScanner is an interface for a *sql.DB or *sql.TX or one of the Squirrel interfaces for the generated scan functions
type RowScanner interface {
	Scan(...interface{}) error
}
