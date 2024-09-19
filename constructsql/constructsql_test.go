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

type mockRows struct {
	rows       []mockRow
	closeErr   error
	iterateErr error

	cursor int
	closed bool
}

func (m *mockRows) Scan(dest ...any) error {
	if m.cursor == 0 {
		return errors.New("mockRows.Scan: before first row, must call Next first")
	}
	if m.cursor > len(m.rows) {
		return errors.New("mockRows.Scan: after last row")
	}

	return m.rows[m.cursor-1].Scan(dest...)
}

func (m *mockRows) Next() bool {
	m.cursor++
	return m.cursor <= len(m.rows)
}

func (m *mockRows) Close() error {
	m.closed = true
	return m.closeErr
}

func (m *mockRows) Err() error {
	return m.iterateErr
}

var _ constructsql.Rows = &mockRows{}

func TestCollectRows(t *testing.T) {
	t.Run("collect rows without error", func(t *testing.T) {
		rows := mockRows{
			rows: []mockRow{
				{
					scanJSON: []byte(`{"id": 1, "name": "test"}`),
				},
			},
		}
		records, err := constructsql.CollectRows[user](&rows, nil)
		require.NoError(t, err)

		assert.Equal(t, []user{{ID: 1, Name: "test"}}, records)
		assert.True(t, rows.closed)
	})

	t.Run("collect empty rows without error", func(t *testing.T) {
		rows := mockRows{
			rows: []mockRow{},
		}
		records, err := constructsql.CollectRows[user](&rows, nil)
		require.NoError(t, err)

		assert.Len(t, records, 0)
		assert.True(t, rows.closed)
	})

	t.Run("collect rows with initial error", func(t *testing.T) {
		rows := mockRows{}
		initialErr := errors.New("some initial error")
		records, err := constructsql.CollectRows[user](&rows, initialErr)
		require.ErrorIs(t, err, initialErr)
		assert.False(t, rows.closed)

		assert.Empty(t, records)
	})

	t.Run("collect rows with scan error", func(t *testing.T) {
		scanErr := errors.New("some scan error")
		rows := mockRows{
			rows: []mockRow{
				{
					scanJSON: []byte(`{"id": 1, "name": "test"}`),
				},
				{
					scanErr: scanErr,
				},
			},
		}
		records, err := constructsql.CollectRows[user](&rows, nil)
		require.ErrorIs(t, err, scanErr)

		assert.Empty(t, records)
		assert.True(t, rows.closed)
	})

	t.Run("collect rows with iterate error", func(t *testing.T) {
		iterateErr := errors.New("some iterate error")
		rows := mockRows{
			rows: []mockRow{
				{
					scanJSON: []byte(`{"id": 1, "name": "test"}`),
				},
			},
			iterateErr: iterateErr,
		}
		records, err := constructsql.CollectRows[user](&rows, nil)
		require.ErrorIs(t, err, iterateErr)

		assert.Empty(t, records)
		assert.True(t, rows.closed)
	})

	t.Run("collect rows with close error", func(t *testing.T) {
		closeErr := errors.New("some error on close")
		rows := mockRows{
			rows: []mockRow{
				{
					scanJSON: []byte(`{"id": 1, "name": "test"}`),
				},
			},
			closeErr: closeErr,
		}
		_, err := constructsql.CollectRows[user](&rows, nil)
		require.ErrorIs(t, err, closeErr)

		assert.True(t, rows.closed)
	})

	t.Run("collect rows with scan and close error", func(t *testing.T) {
		scanErr := errors.New("some scan error")
		closeErr := errors.New("some error on close")
		rows := mockRows{
			rows: []mockRow{
				{
					scanJSON: []byte(`{"id": 1, "name": "test"}`),
				},
				{
					scanErr: scanErr,
				},
			},
			closeErr: closeErr,
		}
		_, err := constructsql.CollectRows[user](&rows, nil)
		require.ErrorIs(t, err, scanErr)
		require.ErrorIs(t, err, closeErr)

		assert.True(t, rows.closed)
	})
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
