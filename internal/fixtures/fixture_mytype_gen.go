// Code generated by construct, DO NOT EDIT.
package fixtures

import (
	"encoding/json"
	uuid "github.com/gofrs/uuid"
	qrb "github.com/networkteam/qrb"
	builder "github.com/networkteam/qrb/builder"
	fn "github.com/networkteam/qrb/fn"
	"time"
)

var myType = struct {
	builder.Identer
	ID         builder.IdentExp
	Foo        builder.IdentExp
	Bar        builder.IdentExp
	Baz        builder.IdentExp
	LastTime   builder.IdentExp
	LastUpdate builder.IdentExp
	Donuts     builder.IdentExp
}{
	Bar:        qrb.N("my_type.the_bar"),
	Baz:        qrb.N("my_type.baz"),
	Donuts:     qrb.N("my_type.donuts"),
	Foo:        qrb.N("my_type.foo"),
	ID:         qrb.N("my_type.id"),
	Identer:    qrb.N("my_type"),
	LastTime:   qrb.N("my_type.last_time"),
	LastUpdate: qrb.N("my_type.updated_at"),
}

var myTargetTypeSortFields = map[string]builder.IdentExp{
	"foo":        myType.Foo,
	"lasttime":   myType.LastTime,
	"lastupdate": myType.LastUpdate,
}

type MyTargetTypeChangeSet struct {
	ID         *uuid.UUID
	Foo        *string
	Bar        []byte
	Baz        *MyEmbeddedType
	LastTime   **time.Time
	LastUpdate *time.Time
	Donuts     []Donut
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
	if c.LastUpdate != nil {
		m["updated_at"] = *c.LastUpdate
	}
	if c.Donuts != nil {
		data, _ := json.Marshal(c.Donuts)
		m["donuts"] = data
	}
	return m
}

func MyTargetTypeToChangeSet(r MyType) (c MyTargetTypeChangeSet) {
	if r.ID != uuid.Nil {
		c.ID = &r.ID
	}
	c.Foo = &r.Foo
	c.Bar = r.Bar
	c.Baz = &r.Baz
	c.LastTime = &r.LastTime
	if !r.LastUpdate.IsZero() {
		c.LastUpdate = &r.LastUpdate
	}
	c.Donuts = r.Donuts
	return
}

var myTargetTypeDefaultJson = fn.JsonBuildObject().
	Prop("ID", myType.ID).
	Prop("Foo", myType.Foo).
	Prop("Bar", qrb.Func("ENCODE", myType.Bar, qrb.String("BASE64"))).
	Prop("Baz", myType.Baz).
	Prop("LastTime", myType.LastTime).
	Prop("LastUpdate", myType.LastUpdate).
	Prop("Donuts", myType.Donuts)