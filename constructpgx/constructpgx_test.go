package constructpgx_test

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/construct/v2"
	"github.com/networkteam/construct/v2/constructpgx"
)

type user struct {
	ID   int
	Name string
}

type mockRow struct {
	scanUser *user
	scanErr  error
}

func (r *mockRow) Scan(dest ...any) error {
	if r.scanErr != nil {
		return r.scanErr
	}

	if len(dest) != 1 {
		return errors.New("mockRow.Scan: dest must have length 1")
	}

	if r.scanUser == nil {
		return errors.New("mockRow.Scan: scanUser is nil")
	}

	*dest[0].(*user) = *r.scanUser

	return nil
}

func TestScanRow(t *testing.T) {
	t.Run("scans row without error", func(t *testing.T) {
		row := mockRow{
			scanUser: &user{ID: 1, Name: "test"},
		}
		record, err := constructpgx.ScanRow[user](&row, nil)
		require.NoError(t, err)

		assert.Equal(t, user{ID: 1, Name: "test"}, record)
	})

	t.Run("returns QueryRow error", func(t *testing.T) {
		row := mockRow{}
		queryRowErr := errors.New("some query error")
		_, err := constructpgx.ScanRow[user](&row, queryRowErr)
		require.ErrorIs(t, err, queryRowErr)
	})

	t.Run("returns Scan error", func(t *testing.T) {
		scanErr := errors.New("some scan error")
		row := mockRow{
			scanErr: scanErr,
		}
		_, err := constructpgx.ScanRow[user](&row, nil)
		require.ErrorIs(t, err, scanErr)
	})

	t.Run("converts Scan err to ErrNotFound", func(t *testing.T) {
		scanErr := pgx.ErrNoRows
		row := mockRow{
			scanErr: scanErr,
		}
		_, err := constructpgx.ScanRow[user](&row, nil)
		require.ErrorIs(t, err, construct.ErrNotFound)
	})

	t.Run("converts QueryRow err to ErrNotFound", func(t *testing.T) {
		// Note: the underlying error usually only occurs during Scan, but we want to make sure it works here as well.
		queryRowErr := pgx.ErrNoRows
		row := mockRow{}
		_, err := constructpgx.ScanRow[user](&row, queryRowErr)
		require.ErrorIs(t, err, construct.ErrNotFound)
	})
}

func TestAssertRowsAffected(t *testing.T) {
	t.Run("returns nil if rows affected matches expected", func(t *testing.T) {
		err := constructpgx.AssertRowsAffected("update", 2)(pgconn.NewCommandTag("UPDATE 2"))
		require.NoError(t, err)
	})

	t.Run("returns error if rows affected does not match expected", func(t *testing.T) {
		err := constructpgx.AssertRowsAffected("update", 1)(pgconn.NewCommandTag("UPDATE 2"))
		require.Error(t, err)
	})
}
