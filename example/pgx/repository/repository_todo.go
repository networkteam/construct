package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/qrbpgx"

	"github.com/networkteam/construct/v2/example/pgx/model"
)

// todoBuildFindQuery creates a partial builder.SelectBuilder that
// - selects a single JSON result by using todoJson()
// - from the todos table
// - and left join the projects for eagerly fetching the project for each todo in a single query
func todoBuildFindQuery() builder.SelectBuilder {
	return qrb.
		SelectJson(todoJson()).
		From(todo).
		LeftJoin(project).On(todo.projectID.Eq(project.id))
}

// FindTodoByID finds a single todo by id
func FindTodoByID(ctx context.Context, executor qrbpgx.Executor, id uuid.UUID) (result model.Todo, err error) {
	q := todoBuildFindQuery().
		Where(todo.id.Eq(qrb.Arg(id)))

	row, err := qrbpgx.Build(q).WithExecutor(executor).QueryRow(ctx)
	if err != nil {
		return result, err
	}
	return pgxScanRow[model.Todo](row)
}

// FindAllTodos finds all todos sorted by title
func FindAllTodos(ctx context.Context, executor qrbpgx.Executor, filter model.TodosFilter) (result []model.Todo, err error) {
	q := todoBuildFindQuery().
		OrderBy(todo.completedAt).Desc().NullsFirst()

	rows, err := qrbpgx.Build(q).WithExecutor(executor).Query(ctx)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgxCollectRow[model.Todo])
}

// InsertTodo inserts a new todo with values from changeSet
func InsertTodo(ctx context.Context, executor qrbpgx.Executor, changeSet TodoChangeSet) error {
	q := qrb.
		InsertInto(todo).
		SetMap(changeSet.toMap())

	_, err := qrbpgx.Build(q).WithExecutor(executor).Exec(ctx)
	return err
}

// UpdateTodo updates a todo with the given id and changes from changeSet
func UpdateTodo(ctx context.Context, executor qrbpgx.Executor, id uuid.UUID, changeSet TodoChangeSet) error {
	q := qrb.
		Update(todo).
		SetMap(changeSet.toMap()).
		Where(todo.id.Eq(qrb.Arg(id)))

	res, err := qrbpgx.Build(q).WithExecutor(executor).Exec(ctx)
	if err != nil {
		return err
	}

	return assertRowsAffected(res, "update", 1)
}

func todoJson() builder.JsonBuildObjectBuilder {
	// Use the generated default select (JSON object builder) and add another property for the aggregated count
	return todoDefaultJson.
		// This is why it's so nice to select JSON: adding a nested select of all the joined project properties is easy
		Prop("Project", projectDefaultJson)
}
