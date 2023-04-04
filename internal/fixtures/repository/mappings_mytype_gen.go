// Code generated by construct, DO NOT EDIT.
package repository

import (
	"encoding/json"
	uuid "github.com/gofrs/uuid"
	fixtures "github.com/networkteam/construct/v2/internal/fixtures"
	qrb "github.com/networkteam/qrb"
	builder "github.com/networkteam/qrb/builder"
	fn "github.com/networkteam/qrb/fn"
	"time"
)

var (
	myType_id       = qrb.N("my_type.id")
	myType_foo      = qrb.N("my_type.foo")
	myType_bar      = qrb.N("my_type.the_bar")
	myType_baz      = qrb.N("my_type.baz")
	myType_lastTime = qrb.N("my_type.last_time")
)
var myTargetTypeSortFields = map[string]builder.IdentExp{
	"foo":      myType_foo,
	"lasttime": myType_lastTime,
}

type MyTargetTypeChangeSet struct {
	ID       *uuid.UUID
	Foo      *string
	Bar      []byte
	Baz      *fixtures.MyEmbeddedType
	LastTime **time.Time
}

func (c MyTargetTypeChangeSet) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if c.ID != nil {
		m["id"] = *c.ID
	}
	if c.Foo != nil {
		m["foo"] = *c.Foo
	}
	if c.Bar != nil {
		m["the_bar"] = c.Bar
	}
	if c.Baz != nil {
		data, _ := json.Marshal(c.Baz)
		m["baz"] = data
	}
	if c.LastTime != nil {
		m["last_time"] = *c.LastTime
	}
	return m
}

func MyTargetTypeToChangeSet(r fixtures.MyType) (c MyTargetTypeChangeSet) {
	if r.ID != uuid.Nil {
		c.ID = &r.ID
	}
	c.Foo = &r.Foo
	c.Bar = r.Bar
	c.Baz = &r.Baz
	c.LastTime = &r.LastTime
	return
}

var myTargetTypeDefaultJson = fn.JsonBuildObject().
	Prop("ID", myType_id).
	Prop("Foo", myType_foo).
	Prop("Bar", qrb.Func("ENCODE", myType_bar, qrb.String("BASE64"))).
	Prop("Baz", myType_baz).
	Prop("LastTime", myType_lastTime)
