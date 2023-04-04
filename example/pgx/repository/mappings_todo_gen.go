// Code generated by construct, DO NOT EDIT.
package repository

import (
	uuid "github.com/gofrs/uuid"
	model "github.com/networkteam/construct/v2/example/pgx/model"
	qrb "github.com/networkteam/qrb"
	builder "github.com/networkteam/qrb/builder"
	fn "github.com/networkteam/qrb/fn"
	"time"
)

var (
	todo_id          = qrb.N("todos.id")
	todo_projectID   = qrb.N("todos.project_id")
	todo_title       = qrb.N("todos.title")
	todo_completedAt = qrb.N("todos.completed_at")
)
var todoSortFields = map[string]builder.IdentExp{}

type TodoChangeSet struct {
	ID          *uuid.UUID
	ProjectID   *uuid.UUID
	Title       *string
	CompletedAt **time.Time
}

func (c TodoChangeSet) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if c.ID != nil {
		m["id"] = *c.ID
	}
	if c.ProjectID != nil {
		m["project_id"] = *c.ProjectID
	}
	if c.Title != nil {
		m["title"] = *c.Title
	}
	if c.CompletedAt != nil {
		m["completed_at"] = *c.CompletedAt
	}
	return m
}

func TodoToChangeSet(r model.Todo) (c TodoChangeSet) {
	if r.ID != uuid.Nil {
		c.ID = &r.ID
	}
	if r.ProjectID != uuid.Nil {
		c.ProjectID = &r.ProjectID
	}
	c.Title = &r.Title
	c.CompletedAt = &r.CompletedAt
	return
}

var todoDefaultJson = fn.JsonBuildObject().
	Prop("ID", todo_id).
	Prop("ProjectID", todo_projectID).
	Prop("Title", todo_title).
	Prop("CompletedAt", todo_completedAt)
