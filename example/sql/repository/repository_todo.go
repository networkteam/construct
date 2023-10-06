package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/qrbsql"

	"github.com/networkteam/construct/v2/constructsql"
	"github.com/networkteam/construct/v2/example/sql/model"
)

type TodoQueryOpts struct {
	IncludeProject bool
	ProjectOpts    ProjectQueryOpts
}

// todoBuildFindQuery creates a partial builder.SelectBuilder that
// - selects a single JSON result by using todoJson()
// - from the todos table
// - and left join the projects for eagerly fetching the project for each todo in a single query (if opts.IncludeProject is true)
func todoBuildFindQuery(opts TodoQueryOpts) builder.SelectBuilder {
	return qrb.
		SelectJson(todoJson(opts)).
		From(todo).
		ApplyIf(opts.IncludeProject, func(q builder.SelectBuilder) builder.SelectBuilder {
			return q.LeftJoin(project).On(todo.ProjectID.Eq(project.ID))
		})
}

// FindTodoByID finds a single todo by id
func FindTodoByID(ctx context.Context, executor qrbsql.Executor, id uuid.UUID, opts TodoQueryOpts) (result model.Todo, err error) {
	q := todoBuildFindQuery(opts).
		Where(todo.ID.Eq(qrb.Arg(id)))

	return constructsql.ScanRow[model.Todo](
		qrbsql.Build(q).WithExecutor(executor).QueryRow(ctx),
	)
}

// FindAllTodos finds all todos sorted by title
func FindAllTodos(ctx context.Context, executor qrbsql.Executor, filter model.TodosFilter, opts TodoQueryOpts) (result []model.Todo, err error) {
	q := todoBuildFindQuery(opts).
		OrderBy(todo.CompletedAt).Desc().NullsFirst().
		SelectBuilder

	if filter.ProjectID != nil {
		q = q.Where(todo.ProjectID.Eq(qrb.Arg(filter.ProjectID)))
	}

	return constructsql.CollectRows[model.Todo](
		qrbsql.Build(q).WithExecutor(executor).Query(ctx),
	)
}

// InsertTodo inserts a new todo with values from changeSet
func InsertTodo(ctx context.Context, executor qrbsql.Executor, changeSet TodoChangeSet) error {
	q := qrb.
		InsertInto(todo).
		SetMap(changeSet.toMap())

	_, err := qrbsql.Build(q).WithExecutor(executor).Exec(ctx)
	return err
}

// UpdateTodo updates a todo with the given id and changes from changeSet
func UpdateTodo(ctx context.Context, executor qrbsql.Executor, id uuid.UUID, changeSet TodoChangeSet) error {
	q := qrb.
		Update(todo).
		SetMap(changeSet.toMap()).
		Where(todo.ID.Eq(qrb.Arg(id)))

	res, err := qrbsql.Build(q).WithExecutor(executor).Exec(ctx)
	if err != nil {
		return err
	}

	return constructsql.AssertRowsAffected(res, "update", 1)
}

func todoJson(opts TodoQueryOpts) builder.JsonBuildObjectBuilder {
	// Use the generated default select (JSON object builder) and add another property for the aggregated count
	return todoDefaultJson.
		// This is why it's so nice to select JSON: adding a nested select of all the joined project properties is easy
		PropIf(opts.IncludeProject, "Project", projectJson(opts.ProjectOpts))
}
