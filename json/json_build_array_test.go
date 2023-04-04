package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/construct/v2/json"
)

func TestJsonBuildArray_Empty(t *testing.T) {
	b := json.JsonBuildArray()

	sql := b.ToSql()
	assert.Equal(t, `JSON_BUILD_ARRAY()`, sql)
}

func TestJsonBuildArray_Add_Immutable(t *testing.T) {
	b1 := json.JsonBuildArray().
		Append(json.Exp("customers.customer_id"))
	b2 := b1.Append(json.Exp("customers.name"))

	sql1 := b1.ToSql()
	assert.Equal(t, `JSON_BUILD_ARRAY(customers.customer_id)`, sql1)

	sql2 := b2.ToSql()
	assert.Equal(t, `JSON_BUILD_ARRAY(customers.customer_id,customers.name)`, sql2)
}
