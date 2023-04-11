package fixtures

import "github.com/gofrs/uuid"

type MyOtherType struct {
	ID uuid.UUID `read_col:"my_other_type.id" write_col:"id"`
}

type NoStructMappingType struct {
	ID uuid.UUID
}
