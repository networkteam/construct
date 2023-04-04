package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"

	"github.com/networkteam/construct/v2/example/pgx/model"
	"github.com/networkteam/construct/v2/example/pgx/repository"
)

func main() {
	var conn *pgx.Conn

	app := &cli.App{
		Name:  "example-pgx",
		Usage: "construct example with PGX, CLI to manage projects and todos",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "postgres-dsn",
				Value: "postgresql://localhost/example-pgx",
			},
		},
		Before: func(c *cli.Context) error {
			var err error
			conn, err = connectDB(c)
			if err != nil {
				return err
			}

			return nil
		},
		After: func(c *cli.Context) error {
			if conn != nil {
				err := conn.Close(c.Context)
				if err != nil {
					return errors.Wrap(err, "closing DB connection")
				}
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "project",
				Subcommands: []*cli.Command{
					{
						Name: "add",
						Action: func(c *cli.Context) error {
							project := model.Project{
								ID:    uuid.Must(uuid.NewV4()),
								Title: c.Args().Get(0),
							}
							err := repository.InsertProject(c.Context, conn, repository.ProjectToChangeSet(project))
							if err != nil {
								return errors.Wrap(err, "inserting project")
							}

							return nil
						},
					},
					{
						Name: "list",
						Action: func(c *cli.Context) error {
							projects, err := repository.FindAllProjects(c.Context, conn)
							if err != nil {
								return errors.Wrap(err, "finding all projects")
							}

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"ID", "Title", "# Todos"})

							for _, project := range projects {
								table.Append([]string{
									project.ID.String(),
									project.Title,
									strconv.Itoa(project.TodoCount),
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
								return errors.Wrap(err, "parsing id")
							}

							project, err := repository.FindProjectByID(c.Context, conn, id)
							if err != nil {
								return errors.Wrap(err, "finding project")
							}

							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"ID", "Title", "# Todos"})

							table.Append([]string{
								project.ID.String(),
								project.Title,
								strconv.Itoa(project.TodoCount),
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
								return fmt.Errorf("expected exactly 2 arguments")
							}
							projectID, err := uuid.FromString(c.Args().Get(0))
							if err != nil {
								return errors.Wrap(err, "parsing project id")
							}
							title := c.Args().Get(1)

							todo := model.Todo{
								ID:        uuid.Must(uuid.NewV4()),
								ProjectID: projectID,
								Title:     title,
							}
							err = repository.InsertTodo(c.Context, conn, repository.TodoToChangeSet(todo))
							if err != nil {
								return errors.Wrap(err, "inserting todo")
							}

							return nil
						},
					},
					{
						Name: "list",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "project-id",
								Usage: "Show todos for the given project",
							},
						},
						Action: func(c *cli.Context) error {
							var filter model.TodosFilter
							if projectIDStr := c.String("project-id"); projectIDStr != "" {
								projectID, err := uuid.FromString(c.Args().Get(0))
								if err != nil {
									return errors.Wrap(err, "parsing project id")
								}
								filter.ProjectID = &projectID
							}

							todos, err := repository.FindAllTodos(c.Context, conn, filter)
							if err != nil {
								return errors.Wrap(err, "finding all todos")
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
								return fmt.Errorf("expected exactly 1 argument")
							}
							todoID, err := uuid.FromString(c.Args().Get(0))
							if err != nil {
								return errors.Wrap(err, "parsing todo id")
							}

							now := time.Now()
							completedAt := &now
							err = repository.UpdateTodo(c.Context, conn, todoID, repository.TodoChangeSet{
								CompletedAt: &completedAt,
							})
							if err != nil {
								return errors.Wrap(err, "updating todo")
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

func connectDB(c *cli.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(c.Context, c.String("postgres-dsn"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	return conn, nil
}
