package repository

import (
	"context"

	"github.com/gofrs/uuid"
	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/qrbsql"

	"github.com/networkteam/construct/v2/constructsql"
	"github.com/networkteam/construct/v2/example/sql/model"
)

type ProjectQueryOpts struct {
	IncludeTodoCount bool
}

// projectBuildFindQuery creates a partial builder.SelectBuilder that
// - selects a single JSON result by using projectJson (which selects a TodoCount if opts.IncludeTodoCount is true)
// - from the projects table
func projectBuildFindQuery(opts ProjectQueryOpts) builder.SelectBuilder {
	return SelectJson(projectJson(opts)).
		From(project).
		SelectBuilder
}

// FindProjectByID finds a single project by id
func FindProjectByID(ctx context.Context, executor qrbsql.Executor, id uuid.UUID, opts ProjectQueryOpts) (result model.Project, err error) {
	q := projectBuildFindQuery(opts).
		Where(project.ID.Eq(Arg(id)))

	return constructsql.ScanRow[model.Project](
		qrbsql.Build(q).WithExecutor(executor).QueryRow(ctx),
	)
}

// FindAllProjects finds all projects sorted by title
func FindAllProjects(ctx context.Context, executor qrbsql.Executor, opts ProjectQueryOpts) (result []model.Project, err error) {
	q := projectBuildFindQuery(opts).
		OrderBy(project.Title)

	return constructsql.CollectRows[model.Project](
		qrbsql.Build(q).WithExecutor(executor).Query(ctx),
	)
}

// InsertProject inserts a new project from a ProjectChangeSet
func InsertProject(ctx context.Context, executor qrbsql.Executor, changeSet ProjectChangeSet) error {
	q := InsertInto(project).
		SetMap(changeSet.toMap())

	_, err := qrbsql.Build(q).WithExecutor(executor).Exec(ctx)
	return err
}

func projectJson(opts ProjectQueryOpts) builder.JsonBuildObjectBuilder {
	// Use the generated default select (JSON object builder) and add another property for the aggregated count#
	// (if opts.IncludeTodoCount is true).
	return projectDefaultJson.
		PropIf(
			opts.IncludeTodoCount,
			"TodoCount",
			Select(fn.Count(N("*"))).
				From(todo).
				Where(todo.ProjectID.Eq(project.ID)),
		)
}
