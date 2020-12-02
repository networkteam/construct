// Code generated by construct, DO NOT EDIT.
package repository

import (
	"database/sql"
	"encoding/json"
	uuid "github.com/gofrs/uuid"
	construct "github.com/networkteam/construct"
	model "github.com/networkteam/construct/example/pgx/model"
	cjson "github.com/networkteam/construct/json"
	"time"
)

const (
	todo_id          = "todos.id"
	todo_projectID   = "todos.project_id"
	todo_title       = "todos.title"
	todo_completedAt = "todos.completed_at"
)

var todoSortFields = map[string]string{}

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

var todoDefaultSelectJson = cjson.JsonBuildObject().
	Set("ID", cjson.Exp("todos.id")).
	Set("ProjectID", cjson.Exp("todos.project_id")).
	Set("Title", cjson.Exp("todos.title")).
	Set("CompletedAt", cjson.Exp("todos.completed_at"))

func todoScanJsonRow(row construct.RowScanner) (result model.Todo, err error) {
	var data []byte
	if err := row.Scan(&data); err != nil {
		if err == sql.ErrNoRows {
			return result, construct.ErrNotFound
		}
		return result, err
	}
	return result, json.Unmarshal(data, &result)
}