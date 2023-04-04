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

// todoBuildFindQuery creates a partial squirrel.SelectBuilder that
// - selects a single JSON result by using buildTodoJson
// - from the todos table
// - and left join the projects for eagerly fetching the project for each todo in a single query
func todoBuildFindQuery() builder.SelectBuilder {
	return qrb.
		SelectJson(todoJson()).
		From(qrb.N("todos")).
		LeftJoin(qrb.N("projects")).On(todo_projectID.Eq(project_id))
}

// FindTodoByID finds a single todo by id
func FindTodoByID(ctx context.Context, executor qrbpgx.Executor, id uuid.UUID) (result model.Todo, err error) {
	q := todoBuildFindQuery().
		Where(todo_id.Eq(qrb.Arg(id)))

	row, err := qrbpgx.Build(q).WithExecutor(executor).QueryRow(ctx)
	if err != nil {
		return result, err
	}
	return pgxScanRow[model.Todo](row)
}

// FindAllTodos finds all todos sorted by title
func FindAllTodos(ctx context.Context, executor qrbpgx.Executor, filter model.TodosFilter) (result []model.Todo, err error) {
	q := todoBuildFindQuery().
		OrderBy(todo_completedAt).Desc().NullsFirst()

	rows, err := qrbpgx.Build(q).WithExecutor(executor).Query(ctx)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgxCollectRow[model.Todo])
}

// InsertTodo inserts a new todo with values from changeSet
func InsertTodo(ctx context.Context, executor qrbpgx.Executor, changeSet TodoChangeSet) error {
	/*
		q := queryBuilder().
			Insert("todos").
			SetMap(changeSet.toMap())
		_, err := pgxExec(ctx, querier, q)
		return err
	*/

	// TODO Implement insert

	return nil
}

// UpdateTodo updates a todo with the given id and changes from changeSet
func UpdateTodo(ctx context.Context, executor qrbpgx.Executor, id uuid.UUID, changeSet TodoChangeSet) error {
	/*
		q := queryBuilder().
			Update("todos").
			Where(squirrel.Eq{todo_id: id}).
			SetMap(changeSet.toMap())
		res, err := pgxExec(ctx, querier, q)
		if err != nil {
			return err
		}
		return assertRowsAffected(res, "update", 1)
	*/

	// TODO Implement update

	return nil
}

func todoJson() builder.JsonBuildObjectBuilder {
	// Use the generated default select (JSON object builder) and add another property for the aggregated count
	return todoDefaultJson.
		// This is why it's so nice to select JSON: adding a nested select of all the joined project properties is easy
		Prop("Project", projectDefaultJson)
}
