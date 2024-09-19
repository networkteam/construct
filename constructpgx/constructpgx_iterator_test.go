//go:build go1.23

package constructpgx_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/construct/v2/constructpgx"
)

func TestIterateRows(t *testing.T) {
	t.Run("iterate rows without error", func(t *testing.T) {
		rows := mockRows{
			rows: []mockRow{
				{
					scanUser: &user{ID: 1, Name: "test"},
				},
			},
		}

		count := 0
		for record, err := range constructpgx.IterateRows[user](&rows, nil) {
			require.NoError(t, err)
			assert.Equal(t, user{ID: 1, Name: "test"}, record)
			count++
		}
		assert.Equal(t, 1, count)

		assert.True(t, rows.closed)
	})

	t.Run("iterate empty rows without error", func(t *testing.T) {
		rows := mockRows{
			rows: []mockRow{},
		}

		count := 0
		for range constructpgx.IterateRows[user](&rows, nil) {
			count++
		}
		assert.Equal(t, 0, count)

		assert.True(t, rows.closed)
	})

	t.Run("iterate rows with early break", func(t *testing.T) {
		rows := mockRows{
			rows: []mockRow{
				{
					scanUser: &user{ID: 1, Name: "test"},
				},
				{
					scanUser: &user{ID: 2, Name: "test"},
				},
			},
		}

		count := 0
		for record, err := range constructpgx.IterateRows[user](&rows, nil) {
			if count == 0 {
				require.NoError(t, err)
				assert.Equal(t, user{ID: 1, Name: "test"}, record)
			} else {
				break
			}
			count++
		}
		assert.Equal(t, 1, count)

		assert.True(t, rows.closed)
	})

	t.Run("iterate rows with early break and iterate error", func(t *testing.T) {
		iterateErr := errors.New("some iterate error")
		rows := mockRows{
			rows: []mockRow{
				{
					scanUser: &user{ID: 1, Name: "test"},
				},
				{
					scanUser: &user{ID: 2, Name: "test"},
				},
			},
			iterateErr: iterateErr,
		}

		count := 0
		for record, err := range constructpgx.IterateRows[user](&rows, nil) {
			if count == 0 {
				require.NoError(t, err)
				assert.Equal(t, user{ID: 1, Name: "test"}, record)
			} else {
				break
			}
			count++
		}
		assert.Equal(t, 1, count)

		assert.True(t, rows.closed)
	})

	t.Run("iterate rows with initial error", func(t *testing.T) {
		rows := mockRows{}
		initialErr := errors.New("some initial error")
		count := 0
		for _, err := range constructpgx.IterateRows[user](&rows, initialErr) {
			assert.ErrorIs(t, err, initialErr)
			count++
		}
		assert.Equal(t, 1, count)

		assert.False(t, rows.closed)
	})

	t.Run("iterate rows with scan error", func(t *testing.T) {
		scanErr := errors.New("some scan error")
		rows := mockRows{
			rows: []mockRow{
				{
					scanUser: &user{ID: 1, Name: "test"},
				},
				{
					scanErr: scanErr,
				},
			},
		}

		count := 0
		for record, err := range constructpgx.IterateRows[user](&rows, nil) {
			if count == 0 {
				require.NoError(t, err)
				assert.Equal(t, user{ID: 1, Name: "test"}, record)
			} else if count == 1 {
				assert.ErrorIs(t, err, scanErr)
			}
			count++
		}
		assert.Equal(t, 2, count)

		assert.True(t, rows.closed)
	})

	t.Run("iterate rows with iterate error", func(t *testing.T) {
		iterateErr := errors.New("some iterate error")
		rows := mockRows{
			rows: []mockRow{
				{
					scanUser: &user{ID: 1, Name: "test"},
				},
			},
			iterateErr: iterateErr,
		}

		count := 0
		for record, err := range constructpgx.IterateRows[user](&rows, nil) {
			if count == 0 {
				require.NoError(t, err)
				assert.Equal(t, user{ID: 1, Name: "test"}, record)
			} else if count == 1 {
				assert.ErrorIs(t, err, iterateErr)
			}
			count++
		}
		assert.Equal(t, 2, count)

		assert.True(t, rows.closed)
	})
}
