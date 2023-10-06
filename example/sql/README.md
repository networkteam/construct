# Construct example with qrb and database/sql

A very simple todo CLI application that shows how custom repository code can use the generated mappings
from construct. It features eager loading of nested relations and usage of`database/sql`.

## How to run

Create a new PostgreSQL database:

    createdb example-sql

Apply the schema:

    psql example-sql < schema.sql 

Add a project:

    go run . project add neverfinished
    
List projects:

    go run . project list

Add a todo:

    go run . todo add [project id] "Write more documentation"
    
List todos:

    go run . todo list

Mark todo as completed:

    go run . todo complete [todo id]
