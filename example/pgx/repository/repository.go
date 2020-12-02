package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/friendsofgo/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"

	"github.com/networkteam/construct"
)

func queryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar)
}

func assertRowsAffected(res pgconn.CommandTag, op string, numberOfRows int64) error {
	rowsAffected := res.RowsAffected()
	if rowsAffected != numberOfRows {
		return errors.Errorf("%s affected %d rows, but expected exactly %d", op, rowsAffected, numberOfRows)
	}
	return nil
}

func pgxQuery(ctx context.Context, querier pgxtype.Querier, q squirrel.SelectBuilder) (pgx.Rows, error) {
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building SQL")
	}
	rows, err := querier.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying SQL")
	}
	return rows, err
}

func pgxExec(ctx context.Context, querier pgxtype.Querier, q squirrel.Sqlizer) (pgconn.CommandTag, error) {
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building SQL")
	}
	res, err := querier.Exec(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "executing SQL")
	}
	return res, nil
}

func pgxQueryRow(ctx context.Context, querier pgxtype.Querier, q squirrel.Sqlizer) (pgx.Row, error) {
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building SQL")
	}
	row := querier.QueryRow(ctx, sql, args...)
	return row, nil
}

func pgxScanRow(scanner construct.RowScanner, result interface{}) error {
	err := scanner.Scan(result)
	if err != nil {
		return errors.Wrap(pgxToConstructErr(err), "scanning row")
	}
	return nil
}

func pgxToConstructErr(err error) error {
	if err == pgx.ErrNoRows {
		return construct.ErrNotFound
	}
	return err
}
