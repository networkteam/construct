//go:build go1.23

package constructsql

import (
	"errors"
	"iter"
	"log"
)

// IterateRows returns an iterator over the rows of a database query and scans them to T.
// It returns a single-use iterator.
func IterateRows[T any](rows Rows, err error) iter.Seq2[T, error] {
	if err != nil {
		return func(yield func(T, error) bool) {
			var result T
			yield(result, err)
		}
	}

	return func(yield func(T, error) bool) {
		var err error
		var iteratorClosed bool

		defer func() {
			closeErr := rows.Close()
			if iteratorClosed && closeErr != nil {
				log.Printf("constructsql: Error closing rows after function for loop body returned false: %v", closeErr)
				return
			}

			err = errors.Join(err, closeErr)
			if err != nil {
				var result T
				yield(result, err)
			}
		}()

		var value T
		for rows.Next() {
			value, err = scanRow[T](rows)
			if err != nil {
				return
			}
			if !yield(value, nil) {
				iteratorClosed = true
				return
			}
		}

		if err = rows.Err(); err != nil {
			return
		}
	}
}
