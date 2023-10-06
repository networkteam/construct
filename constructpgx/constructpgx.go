package constructpgx

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/networkteam/construct/v2"
)

// CollectRows collects all rows to the given target type from a ExecutiveQueryBuilder.Query result.
// It is based on pgx.CollectRows but accepts an additional error for easier use.
func CollectRows[T any](rows pgx.Rows, err error) ([]T, error) {
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

func pgxToConstructErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return construct.ErrNotFound
	}
	return err
}

// ScanRow scans a single row to the given target type from a ExecutiveQueryBuilder.QueryRow result.
func ScanRow[T any](row pgx.Row, err error) (T, error) {
	var result T
	if err != nil {
		return result, pgxToConstructErr(err)
	}
	return scanRow[T](row)
}

func scanRow[T any](row pgx.Row) (T, error) {
	var result T
	err := row.Scan(&result)
	if err != nil {
		return result, fmt.Errorf("scanning row: %w", pgxToConstructErr(err))
	}
	return result, nil
}

// AssertRowsAffected checks if the given result affected exactly the expected number of rows.
func AssertRowsAffected(result pgconn.CommandTag, operation string, expectedRows int) error {
	actualRows := result.RowsAffected()
	if actualRows != int64(expectedRows) {
		return fmt.Errorf("%s affected %d rows, but expected exactly %d", operation, actualRows, expectedRows)
	}

	return nil

}
