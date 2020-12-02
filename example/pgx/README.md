# Construct example with Squirrel and pgx

A very simple todo CLI application that shows how custom repository code can use the generated mappings
from construct. It features eager loading of nested relations and usage of the pgx features instead of
`database/sql` (which can be used in the same way with construct).

## How to run

Create a new PostgreSQL database:

    createdb example-pgx 

Apply the schema:

    psql example-pgx < schema.sql 

Add a project:

    go run . project add neverfinished
    
List projects:

    go run . project list

Add a todo:

    go run . todo add [project id] "Write more documentation"
    
List todos:

    go run . todo list
