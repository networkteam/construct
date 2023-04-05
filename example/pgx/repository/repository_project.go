package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/qrbpgx"

	"github.com/networkteam/construct/v2/example/pgx/model"
)

// projectBuildFindQuery creates a partial builder.SelectBuilder that
// - selects a single JSON result by using buildProjectJson
// - from the projects table
// - and left joins an aggregation of todo counts by project
func projectBuildFindQuery() builder.SelectBuilder {
	return qrb.
		SelectJson(projectJson()).
		From(qrb.N("projects")).
		LeftJoin(
			qrb.Select(fn.Count(todo_projectID)).As("count").
				Select(todo_projectID).
				From(qrb.N("todos")).
				GroupBy(todo_projectID),
		).As("todo_counts").On(project_id.Eq(qrb.N("todo_counts.project_id")))
}

// FindProjectByID finds a single project by id
func FindProjectByID(ctx context.Context, executor qrbpgx.Executor, id uuid.UUID) (result model.Project, err error) {
	q := projectBuildFindQuery().
		Where(project_id.Eq(qrb.Arg(id)))

	row, err := qrbpgx.Build(q).WithExecutor(executor).QueryRow(ctx)
	if err != nil {
		return result, err
	}
	return pgxScanRow[model.Project](row)
}

// FindAllProjects finds all projects sorted by title
func FindAllProjects(ctx context.Context, executor qrbpgx.Executor) (result []model.Project, err error) {
	q := projectBuildFindQuery().
		OrderBy(project_title)

	rows, err := qrbpgx.Build(q).WithExecutor(executor).Query(ctx)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgxCollectRow[model.Project])
}

// InsertProject inserts a new project from a ProjectChangeSet
func InsertProject(ctx context.Context, executor qrbpgx.Executor, changeSet ProjectChangeSet) error {
	q := qrb.
		InsertInto("projects").
		SetMap(changeSet.toMap())

	_, err := qrbpgx.Build(q).WithExecutor(executor).Exec(ctx)
	return err
}

func projectJson() builder.JsonBuildObjectBuilder {
	// Use the generated default select (JSON object builder) and add another property for the aggregated count
	return projectDefaultJson.
		// Wrap with COALESCE for null values because of LEFT JOIN (if no todos are present for a project)
		Prop("TodoCount", qrb.Coalesce(qrb.N("todo_counts.count"), qrb.Int(0)))
}
