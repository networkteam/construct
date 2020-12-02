package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/construct/json"
)

func TestExpIfNotNull(t *testing.T) {
	g := json.ExpIfNotNull("related_id", json.JsonBuildObject().Set("ID", json.Exp("related_id")))

	sql := json.GeneratorToSql(g)
	assert.Equal(t, "CASE WHEN related_id IS NOT NULL THEN JSON_BUILD_OBJECT('ID',related_id) END", sql)
}
