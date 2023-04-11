package model

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/networkteam/construct/v2"
)

// Todo model belonging to a project
type Todo struct {
	construct.Table `table_name:"todos"` // The table_name struct tag could also be applied to an empty embedded struct
	ID              uuid.UUID            `read_col:"todos.id" write_col:"id"`
	ProjectID       uuid.UUID            `read_col:"todos.project_id" write_col:"project_id"`
	Title           string               `read_col:"todos.title" write_col:"title"`
	CompletedAt     *time.Time           `read_col:"todos.completed_at" write_col:"completed_at"`

	// for eager loading of referenced records
	Project *Project
}

// TodosFilter for filtering todos
type TodosFilter struct {
	ProjectID *uuid.UUID
}
