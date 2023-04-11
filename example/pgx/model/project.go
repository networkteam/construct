package model

import (
	"github.com/gofrs/uuid"
)

// Project model with todos
type Project struct {
	ID    uuid.UUID `table_name:"projects" read_col:"projects.id" write_col:"id"`
	Title string    `read_col:"projects.title" write_col:"title"`

	// TodoCount will be set to the count of todos for this project (if ProjectQueryOpts.IncludeTodoCount is set)
	TodoCount *int
}
