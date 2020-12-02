package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype/pgxtype"

	"github.com/networkteam/construct/example/pgx/model"
	cjson "github.com/networkteam/construct/json"
)

// projectBuildFindQuery creates a partial squirrel.SelectBuilder that
// - selects a single JSON result by using buildProjectJson
// - from the projects table
// - and left joins an aggregation of todo counts by project
func projectBuildFindQuery() squirrel.SelectBuilder {
	return queryBuilder().
		Select(buildProjectJson()).
		From("projects").
		LeftJoin(`(
			SELECT COUNT(todos.project_id) AS count, todos.project_id FROM todos GROUP BY todos.project_id
		) AS todo_counts ON projects.id = todo_counts.project_id`)
}

// FindProjectByID finds a single project by id
func FindProjectByID(ctx context.Context, querier pgxtype.Querier, id uuid.UUID) (result model.Project, err error) {
	q := projectBuildFindQuery().
		Where(squirrel.Eq{project_id: id})

	row, err := pgxQueryRow(ctx, querier, q)
	if err != nil {
		return result, err
	}
	return result, pgxScanRow(row, &result)
}

// FindAllProjects finds all projects sorted by title
func FindAllProjects(ctx context.Context, querier pgxtype.Querier) (result []model.Project, err error) {
	q := projectBuildFindQuery().
		OrderBy(project_title)

	rows, err := pgxQuery(ctx, querier, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var record model.Project
	for rows.Next() {
		err := pgxScanRow(rows, &record)
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return
}

// InsertProject inserts a new project from a ProjectChangeSet
func InsertProject(ctx context.Context, querier pgxtype.Querier, changeSet ProjectChangeSet) error {
	q := queryBuilder().
		Insert("projects").
		SetMap(changeSet.toMap())
	_, err := pgxExec(ctx, querier, q)
	return err
}

func buildProjectJson() string {
	// Use the generated default select (JSON object builder) and add another property for the aggregated count
	return projectDefaultSelectJson.
		// Wrap with COALESCE for null values because of LEFT JOIN (if no todos are present for a project)
		Set("TodoCount", cjson.Exp("COALESCE(todo_counts.count, 0)")).
		ToSql()
}
