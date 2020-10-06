package construct

import "errors"

var ErrNotFound = errors.New("record not found")

type RowScanner interface {
	Scan(...interface{}) error
}
