package constructsql_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/construct/v2"
	"github.com/networkteam/construct/v2/constructsql"
)

type user struct {
	ID   int
	Name string
}

type mockRow struct {
	scanJSON []byte
	scanErr  error
}

func (r *mockRow) Scan(dest ...any) error {
	if r.scanErr != nil {
		return r.scanErr
	}

	if len(dest) != 1 {
		return errors.New("mockRow.Scan: dest must have length 1")
	}

	if r.scanJSON == nil {
		return errors.New("mockRow.Scan: scanJSON is nil")
	}

	*dest[0].(*[]byte) = r.scanJSON

	return nil
}

func TestScanRow(t *testing.T) {
	t.Run("scans row without error", func(t *testing.T) {
		row := mockRow{
			scanJSON: []byte(`{"id": 1, "name": "test"}`),
		}
		record, err := constructsql.ScanRow[user](&row, nil)
		require.NoError(t, err)

		assert.Equal(t, user{ID: 1, Name: "test"}, record)
	})

	t.Run("returns QueryRow error", func(t *testing.T) {
		row := mockRow{}
		queryRowErr := errors.New("some query error")
		_, err := constructsql.ScanRow[user](&row, queryRowErr)
		require.ErrorIs(t, err, queryRowErr)
	})

	t.Run("returns Scan error", func(t *testing.T) {
		scanErr := errors.New("some scan error")
		row := mockRow{
			scanErr: scanErr,
		}
		_, err := constructsql.ScanRow[user](&row, nil)
		require.ErrorIs(t, err, scanErr)
	})

	t.Run("converts Scan err to ErrNotFound", func(t *testing.T) {
		scanErr := sql.ErrNoRows
		row := mockRow{
			scanErr: scanErr,
		}
		_, err := constructsql.ScanRow[user](&row, nil)
		require.ErrorIs(t, err, construct.ErrNotFound)
	})

	t.Run("converts QueryRow err to ErrNotFound", func(t *testing.T) {
		// Note: the underlying error usually only occurs during Scan, but we want to make sure it works here as well.
		queryRowErr := sql.ErrNoRows
		row := mockRow{}
		_, err := constructsql.ScanRow[user](&row, queryRowErr)
		require.ErrorIs(t, err, construct.ErrNotFound)
	})
}

type sqlResult struct {
	affectedRows    int64
	affectedRowsErr error
}

func (r sqlResult) LastInsertId() (int64, error) {
	panic("implement me")
}

func (r sqlResult) RowsAffected() (int64, error) {
	return r.affectedRows, r.affectedRowsErr
}

func TestAssertRowsAffected(t *testing.T) {
	t.Run("returns nil if rows affected matches expected", func(t *testing.T) {
		err := constructsql.AssertRowsAffected("update", 2)(sqlResult{affectedRows: 2}, nil)
		require.NoError(t, err)
	})

	t.Run("returns error if rows affected does not match expected", func(t *testing.T) {
		err := constructsql.AssertRowsAffected("update", 1)(sqlResult{affectedRows: 2}, nil)
		require.Error(t, err)
	})

	t.Run("returns error if rows affected returns error", func(t *testing.T) {
		err := constructsql.AssertRowsAffected("update", 1)(sqlResult{affectedRowsErr: errors.New("some error")}, nil)
		require.Error(t, err)
	})

	t.Run("returns previous error", func(t *testing.T) {
		err := constructsql.AssertRowsAffected("update", 1)(sqlResult{}, errors.New("some error"))
		require.Error(t, err)
	})
}
