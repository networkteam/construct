package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/networkteam/construct/v2/example/sql/model"
	"github.com/networkteam/construct/v2/example/sql/repository"
)

func main() {
	var db *sql.DB

	app := &cli.App{
		Name:  "example-sql",
		Usage: "construct example with database/sql, CLI to manage projects and todos",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "postgres-dsn",
				Value: "postgresql://localhost/example-sql",
			},
		},
		Before: func(c *cli.Context) error {
			var err error
			db, err = connectDB(c)
			if err != nil {
				return err
			}

			return nil
		},
		After: func(c *cli.Context) error {
			if db != nil {
				err := db.Close()
				if err != nil {
					return fmt.Errorf("closing DB connection: %w", err)
				}
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "project",
				Subcommands: []*cli.Command{
					{
						Name:      "add",
						ArgsUsage: "[title]",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return fmt.Errorf("expected exactly 1 argument: [title]")
							}
							project := model.Project{
								ID:    uuid.Must(uuid.NewV4()),
								Title: c.Args().Get(0),
							}
							err := repository.InsertProject(c.Context, db, repository.ProjectToChangeSet(project))
							if err != nil {
								return fmt.Errorf("inserting project: %w", err)
							}

							return nil
						},
					},
					{
						Name: "list",
						Action: func(c *cli.Context) error {
							projects, err := repository.FindAllProjects(c.Context, db, repository.ProjectQueryOpts{
								IncludeTodoCount: true,
							})
							if err != nil {
								return fmt.Errorf("finding all projects: %w", err)
							}

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"ID", "Title", "# Todos"})

							for _, project := range projects {
								table.Append([]string{
									project.ID.String(),
									project.Title,
									strconv.Itoa(*project.TodoCount),
								})
							}

							table.Render()

							return nil
						},
					},
					{
						Name:      "show",
						ArgsUsage: "[id]",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return fmt.Errorf("expected exactly 1 argument: [id]")
							}
							id, err := uuid.FromString(c.Args().Get(0))
							if err != nil {
								return fmt.Errorf("parsing id: %w", err)
							}

							project, err := repository.FindProjectByID(c.Context, db, id, repository.ProjectQueryOpts{
								IncludeTodoCount: true,
							})
							if err != nil {
								return fmt.Errorf("finding project: %w", err)
							}

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"ID", "Title", "# Todos"})

							table.Append([]string{
								project.ID.String(),
								project.Title,
								strconv.Itoa(*project.TodoCount),
							})

							table.Render()

							return nil
						},
					},
				},
			},
			{
				Name: "todo",
				Subcommands: []*cli.Command{
					{
						Name:      "add",
						ArgsUsage: "[project id] [title]",
						Action: func(c *cli.Context) error {
							if c.NArg() != 2 {
								return fmt.Errorf("expected exactly 2 arguments: [project id] [title]")
							}
							projectID, err := uuid.FromString(c.Args().Get(0))
							if err != nil {
								return fmt.Errorf("parsing project id: %w", err)
							}
							title := c.Args().Get(1)

							todo := model.Todo{
								ID:        uuid.Must(uuid.NewV4()),
								ProjectID: projectID,
								Title:     title,
							}
							err = repository.InsertTodo(c.Context, db, repository.TodoToChangeSet(todo))
							if err != nil {
								return fmt.Errorf("inserting todo: %w", err)
							}

							return nil
						},
					},
					{
						Name: "list",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "project-id",
								Usage: "Filter todos by the given project (optional)",
							},
						},
						Action: func(c *cli.Context) error {
							var filter model.TodosFilter
							if projectIDStr := c.String("project-id"); projectIDStr != "" {
								projectID, err := uuid.FromString(projectIDStr)
								if err != nil {
									return fmt.Errorf("parsing project id: %w", err)
								}
								filter.ProjectID = &projectID
							}

							todos, err := repository.FindAllTodos(c.Context, db, filter, repository.TodoQueryOpts{
								IncludeProject: true,
							})
							if err != nil {
								return fmt.Errorf("finding all todos: %w", err)
							}

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"ID", "Project", "Title", "Completed at"})

							for _, todo := range todos {
								completedAt := ""
								if todo.CompletedAt != nil {
									completedAt = todo.CompletedAt.Format(time.RFC822)
								}

								table.Append([]string{
									todo.ID.String(),
									todo.Project.Title,
									todo.Title,
									completedAt,
								})
							}

							table.Render()

							return nil
						},
					},
					{
						Name:      "complete",
						ArgsUsage: "[todo id]",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return fmt.Errorf("expected exactly 1 argument: [todo id]")
							}
							todoID, err := uuid.FromString(c.Args().Get(0))
							if err != nil {
								return fmt.Errorf("parsing todo id: %w", err)
							}

							now := time.Now()
							completedAt := &now
							err = repository.UpdateTodo(c.Context, db, todoID, repository.TodoChangeSet{
								CompletedAt: &completedAt,
							})
							if err != nil {
								return fmt.Errorf("updating todo: %w", err)
							}

							return nil
						},
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func connectDB(c *cli.Context) (*sql.DB, error) {
	db, err := sql.Open("pgx", c.String("postgres-dsn"))
	if err != nil {
		return nil, fmt.Errorf("opening DB connection: %w", err)
	}

	return db, nil
}
