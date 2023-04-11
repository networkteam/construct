package model

import (
	"github.com/gofrs/uuid"
)

// Project model with todos
type Project struct {
	ID    uuid.UUID `table_name:"projects" read_col:"projects.id" write_col:"id"`
	Title string    `read_col:"projects.title" write_col:"title"`

	// for embedded loading of counts
	TodoCount int
}
