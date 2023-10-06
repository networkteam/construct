package constructsql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/networkteam/construct/v2"
)

// CollectRows collects all rows to the given target type from a ExecutiveQueryBuilder.Query result.
func CollectRows[T any](rows Rows, err error) ([]T, error) {
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	slice := []T{}

	for rows.Next() {
		value, err := scanRow[T](rows)
		if err != nil {
			return nil, err
		}
		slice = append(slice, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return slice, nil
}

type Rows interface {
	RowScanner
	Next() bool
	Close() error
	Err() error
}

type RowScanner interface {
	Scan(dest ...any) error
}

func sqlToConstructErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return construct.ErrNotFound
	}
	return err
}

// ScanRow scans a single row to the given target type from a ExecutiveQueryBuilder.QueryRow result.
// It expects a single JSON column to be selected and unmarshals to the given struct type.
func ScanRow[T any](row RowScanner, err error) (T, error) {
	var result T
	if err != nil {
		return result, sqlToConstructErr(err)
	}
	return scanRow[T](row)
}

func scanRow[T any](row RowScanner) (result T, err error) {
	var data []byte
	if err := row.Scan(&data); err != nil {
		return result, fmt.Errorf("scanning row: %w", sqlToConstructErr(err))
	}
	return result, json.Unmarshal(data, &result)
}

// AssertRowsAffected checks if the given result affected exactly the expected number of rows.
func AssertRowsAffected(operation string, expectedRows int) func(sql.Result, error) error {
	return func(result sql.Result, err error) error {
		if err != nil {
			return err
		}
		actualRows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("getting affected rows: %w", err)
		}

		if actualRows != int64(expectedRows) {
			return fmt.Errorf("%s affected %d rows, but expected exactly %d", operation, actualRows, expectedRows)
		}

		return nil
	}
}
