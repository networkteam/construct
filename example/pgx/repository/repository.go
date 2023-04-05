package repository

import (
	"github.com/friendsofgo/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/networkteam/construct/v2"
)

func assertRowsAffected(res pgconn.CommandTag, op string, numberOfRows int64) error {
	rowsAffected := res.RowsAffected()
	if rowsAffected != numberOfRows {
		return errors.Errorf("%s affected %d rows, but expected exactly %d", op, rowsAffected, numberOfRows)
	}
	return nil
}

func pgxCollectRow[T any](row pgx.CollectableRow) (T, error) {
	return pgxScanRow[T](row)
}

func pgxScanRow[T any](row pgx.Row) (T, error) {
	var result T
	err := row.Scan(&result)
	if err != nil {
		return result, errors.Wrap(pgxToConstructErr(err), "scanning row")
	}
	return result, nil
}

func pgxToConstructErr(err error) error {
	if err == pgx.ErrNoRows {
		return construct.ErrNotFound
	}
	return err
}
