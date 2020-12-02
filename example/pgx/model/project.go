package model

import (
	"github.com/gofrs/uuid"
)

type Project struct {
	ID    uuid.UUID `read_col:"projects.id" write_col:"id"`
	Title string    `read_col:"projects.title" write_col:"title"`

	// for embedded loading of counts
	TodoCount int
}
