package fixtures

import (
	"time"

	"github.com/gofrs/uuid"
)

type MyEmbeddedType int

type Donut int

type MyType struct {
	ID         uuid.UUID
	Foo        string
	Bar        []byte
	Baz        MyEmbeddedType
	LastTime   *time.Time
	LastUpdate time.Time
	Donuts     []Donut
}
