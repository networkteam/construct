package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/construct/json"
)

func TestExp(t *testing.T) {
	b := json.Exp("foo.bar")
	sql := json.GeneratorToSql(b)

	assert.Equal(t, `foo.bar`, sql)
}
