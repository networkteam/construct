package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/construct/json"
)

func TestJsonBuildObject_Set_Single(t *testing.T) {
	b := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id"))

	sql := b.ToSql()
	assert.Equal(t, `JSON_BUILD_OBJECT('id',customers.customer_id)`, sql)
}

func TestJsonBuildObject_Set_Multiple(t *testing.T) {
	b := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id")).
		Set("name", json.Exp("customers.name"))

	sql := b.ToSql()
	assert.Equal(t, `JSON_BUILD_OBJECT('id',customers.customer_id,'name',customers.name)`, sql)
}

func TestJsonBuildObject_Set_Override(t *testing.T) {
	b := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id")).
		Set("name", json.Exp("customers.name")).
		Set("name", json.Exp("customers.my_name"))

	sql := b.ToSql()
	assert.Equal(t, `JSON_BUILD_OBJECT('id',customers.customer_id,'name',customers.my_name)`, sql)
}

func TestJsonBuildObject_Set_Immutable(t *testing.T) {
	b1 := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id"))
	b2 := b1.Set("name", json.Exp("customers.name"))

	sql1 := b1.ToSql()
	assert.Equal(t, `JSON_BUILD_OBJECT('id',customers.customer_id)`, sql1)

	sql2 := b2.ToSql()
	assert.Equal(t, `JSON_BUILD_OBJECT('id',customers.customer_id,'name',customers.name)`, sql2)
}

func TestJsonBuildObject_Delete(t *testing.T) {
	b := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id")).
		Set("name", json.Exp("customers.name"))

	b = b.Delete("id")

	sql := b.ToSql()
	assert.Equal(t, `JSON_BUILD_OBJECT('name',customers.name)`, sql)
}

func TestJsonBuildObject_Delete_Immutable(t *testing.T) {
	b1 := json.JsonBuildObject().
		Set("name", json.Exp("customers.name"))

	b2 := b1.Delete("name")

	sql1 := b1.ToSql()
	assert.Equal(t, `JSON_BUILD_OBJECT('name',customers.name)`, sql1)

	sql2 := b2.ToSql()
	assert.Equal(t, `JSON_BUILD_OBJECT()`, sql2)
}
