package json_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/construct/json"
)

func TestJsonBuildObject_Add_Single(t *testing.T) {
	b := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id"))
	sql := b.ToSql()

	assert.Equal(t, `JSON_BUILD_OBJECT('id',customers.customer_id)`, sql)
}

func TestJsonBuildObject_Add_Multiple(t *testing.T) {
	b := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id")).
		Set("name", json.Exp("customers.name"))
	sql := b.ToSql()

	assert.Equal(t, `JSON_BUILD_OBJECT('id',customers.customer_id,'name',customers.name)`, sql)
}

func TestJsonBuildObject_Add_Override(t *testing.T) {
	b := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id")).
		Set("name", json.Exp("customers.name")).
		Set("name", json.Exp("customers.my_name"))
	sql := b.ToSql()

	assert.Equal(t, `JSON_BUILD_OBJECT('id',customers.customer_id,'name',customers.my_name)`, sql)
}

func TestJsonBuildObject_Delete(t *testing.T) {
	b := json.JsonBuildObject().
		Set("id", json.Exp("customers.customer_id")).
		Set("name", json.Exp("customers.name"))
	b.Delete("id")

	sql := b.ToSql()

	assert.Equal(t, `JSON_BUILD_OBJECT('name',customers.name)`, sql)
}
