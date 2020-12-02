package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype/pgxtype"

	"github.com/networkteam/construct/example/pgx/model"
)

// todoBuildFindQuery creates a partial squirrel.SelectBuilder that
// - selects a single JSON result by using buildTodoJson
// - from the todos table
// - and left join the projects for eagerly fetching the project for each todo in a single query
func todoBuildFindQuery() squirrel.SelectBuilder {
	return queryBuilder().
		Select(buildTodoJson()).
		From("todos").
		LeftJoin("projects ON todos.project_id = projects.id")
}

// FindTodoByID finds a single todo by id
func FindTodoByID(ctx context.Context, querier pgxtype.Querier, id uuid.UUID) (result model.Todo, err error) {
	q := todoBuildFindQuery().
		Where(squirrel.Eq{todo_id: id})

	row, err := pgxQueryRow(ctx, querier, q)
	if err != nil {
		return result, err
	}
	return result, pgxScanRow(row, &result)
}

// FindAllTodos finds all todos sorted by title
func FindAllTodos(ctx context.Context, querier pgxtype.Querier, filter model.TodosFilter) (result []model.Todo, err error) {
	q := todoBuildFindQuery().
		OrderBy(fmt.Sprintf("%s DESC NULLS FIRST", todo_completedAt))

	rows, err := pgxQuery(ctx, querier, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var record model.Todo
	for rows.Next() {
		err := pgxScanRow(rows, &record)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return
}

// InsertTodo inserts a new todo with values from changeSet
func InsertTodo(ctx context.Context, querier pgxtype.Querier, changeSet TodoChangeSet) error {
	q := queryBuilder().
		Insert("todos").
		SetMap(changeSet.toMap())
	_, err := pgxExec(ctx, querier, q)
	return err
}

// UpdateTodo updates a todo with the given id and changes from changeSet
func UpdateTodo(ctx context.Context, querier pgxtype.Querier, id uuid.UUID, changeSet TodoChangeSet) error {
	q := queryBuilder().
		Update("todos").
		Where(squirrel.Eq{todo_id: id}).
		SetMap(changeSet.toMap())
	res, err := pgxExec(ctx, querier, q)
	if err != nil {
		return err
	}
	return assertRowsAffected(res, "update", 1)
}

func buildTodoJson() string {
	// Use the generated default select (JSON object builder) and add another property for the aggregated count
	return todoDefaultSelectJson.
		// This is why it's so nice to select JSON: adding a nested select of all the joined project properties is easy
		Set("Project", projectDefaultSelectJson).
		ToSql()
}
