package construct

import "errors"

// ErrNotFound will be returned from generated scan functions if the query result was empty, but a row was scanned.
var ErrNotFound = errors.New("record not found")

// Table is an empty struct for embedding into models to derive the table name from struct tags.
type Table struct{}
