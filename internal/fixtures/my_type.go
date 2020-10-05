package fixtures

import (
	"time"

	"github.com/gofrs/uuid"
)

type MyType struct {
	// ID is readable and writable
	ID uuid.UUID `read_col:"my_type.id" write_col:"id"`
	// Foo is readable and sortable
	Foo string `read_col:"my_type.foo,sortable" write_col:"foo"`
	// Bar is readable and writable (from a column with different name)
	Bar []byte `read_col:"my_type.the_bar" write_col:"the_bar"`
	// LastTime is a readable, sortable and writable pointer column
	LastTime *time.Time `read_col:"my_type.last_time,sortable" write_col:"last_time"`
}
