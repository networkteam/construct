package fixtures

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/networkteam/construct/v2"
)

// MyType is a fixture struct type
type MyType struct {
	construct.Table `table_name:"my_type"`

	// ID is readable and writable
	ID uuid.UUID `read_col:"my_type.id" write_col:"id"`
	// Foo is readable and sortable
	Foo string `read_col:"my_type.foo,sortable" write_col:"foo"`
	// Bar is readable and writable (from a column with different name)
	Bar []byte `read_col:"my_type.the_bar" write_col:"the_bar"`
	// Baz should be embedded as JSON
	Baz MyEmbeddedType `read_col:"my_type.baz" write_col:"baz,json"`
	// LastTime is a readable, sortable and writable pointer column
	LastTime *time.Time `read_col:"my_type.last_time,sortable" write_col:"last_time"`
	// LastUpdate is a readable, sortable and writable non-pointer time.Time column that is also named differently
	LastUpdate time.Time `read_col:"my_type.updated_at,sortable" write_col:"updated_at"`
	// Donuts is a readable and writable slice of non-pointer Donut structs
	Donuts []Donut `read_col:"my_type.donuts" write_col:"donuts,json"`
}

// MyEmbeddedType will be embedded in MyType
type MyEmbeddedType struct {
	Fizz bool
	Buzz bool
}

type Donut struct {
	Flavor string
	Size   int
}
