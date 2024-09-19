//go:build go1.23

package constructpgx

import (
	"errors"
	"iter"
	"log"

	"github.com/jackc/pgx/v5"
)

// IterateRows returns an iterator over the rows of a database query and scans them to T.
// If an error occurred initially, while scanning or closing, it returns a single-use iterator with the error.
// It returns a single-use iterator.
func IterateRows[T any](rows pgx.Rows, err error) iter.Seq2[T, error] {
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
			rows.Close()
			closeErr := rows.Err() // pgx states that rows.Err() should be called after rows.Close() or Next() is false
			err = errors.Join(err, closeErr)
			if err != nil {
				if iteratorClosed {
					log.Printf("constructsql: Unhandled error after closing rows and function for loop body returned false: %v", err)
					return
				}
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
	}
}
