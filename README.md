# construct : Go generators for low abstraction persistence with PostgreSQL

[![Go Reference](https://pkg.go.dev/badge/github.com/networkteam/construct/v2.svg)](https://pkg.go.dev/github.com/networkteam/construct/v2)
[![Build status](https://github.com/networkteam/construct/actions/workflows/test.yml/badge.svg?branch=v2)](https://github.com/networkteam/construct/actions/workflows/test.yml)
[![Coverage](https://codecov.io/gh/networkteam/construct/branch/v2/graph/badge.svg?token=Y0GHTB40GG)](https://codecov.io/gh/networkteam/construct)
[![Go Report Card](https://goreportcard.com/badge/github.com/networkteam/construct/v2)](https://goreportcard.com/report/github.com/networkteam/construct/v2)

## Overview

Got tired of too many abstractions over all the features PostgreSQL provides when using an ORM? But rolling your own persistence code is tedious and there's too much boilderplate?
This is a code generator to generate a bunch of structs and functions to implement persistence code with a few line and keep all the power PostgreSQL provides.

## Example

**Your custom type defines fields for reading and writing**

*model/customer.go*

```go
package model

import "github.com/networkteam/construct/v2"

// Customer is just a struct, only struct tags are used in construct to generate code (no struct embedding needed) 
type Customer struct {
	ID   uuid.UUID `table_name:"customers" read_col:"customers.customer_id" write_col:"customer_id"` // The table_name tag is optional and can be specified at most once per struct
	Name string    `read_col:"customers.name,sortable" write_col:"name"`
	// Fields can be serialized as JSON by adding a "json" option to the "write_col" tag.
	// It works perfectly with a column of type jsonb or json in PostgreSQL.
	ContactPerson Contact `read_col:"customers.contact_person,sortable" write_col:"contact_person,json"`

	// ProjectCount is not mapped to the table but used in the select for reading an aggregate count
	ProjectCount int

	CreatedAt time.Time `read_col:"customers.created_at,sortable" write_col:"created_at"`
	UpdatedAt time.Time `read_col:"customers.updated_at,sortable"`
}

// Contact is an embedded type and doesn't need any tags
type Contact struct {
	FirstName  string
	MiddleName string
	LastName   string
}
```

**Construct generates structs and functions that help to read and insert / update data**

Generate code in your persistence package:

*repository/mappings.go*
```go
package repository

//go:generate go run github.com/networkteam/construct/v2/cmd/construct my/project/model
```

```bash
go generate ./repository
```

**Roll your own persistence code for full control and low abstraction**

Structure your persistence code as it fits the project. A very simple and working approach is to have a bunch of
functions that operator on the target type for finding, inserting, updating and deleting. Add more complex queries
as you like in the same way. 

Construct will automatically generate identifier expressions of fields (read columns) via [github.com/networkteam/qrb](https://github.com/networkteam/qrb)
and a default `json_build_object` expression to select a model via JSON. This can be further modified to add additional properties. 

*repository/customer_repository.go*
```go
package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/qrbpgx"

	".../myproject/model"
)

func FindCustomerByID(ctx context.Context, executor qrbpgx.Executor, id uuid.UUID) (model.Customer, error) {
	q := qrb.
		Select(customerJson()).
		// A schema type for read columns (and the table) is generated as qrb identifier expressions by construct
		From(customer).
		LeftJoin(project).On(project.CustomerID.Eq(customer.ID)).
		Where(customer.ID.Eq(qrb.Arg(id))).
		GroupBy(customer.ID)

	row, err := qrbpgx.Build(q).WithExecutor(executor).QueryRow(ctx)
	if err != nil {
		return result, err
	}
	return pgxScanRow[model.Customer](row)
}

// CustomerChangeSet is generated by construct for handling partially filled models
func InsertCustomer(ctx context.Context, executor qrbpgx.Executor, changeSet CustomerChangeSet) error {
	q := qrb.
		InsertInto(customer).
		// toMap is generated by construct 
		SetMap(changeSet.toMap())

	_, err := qrbpgx.Build(q).WithExecutor(executor).Exec(ctx)
	return err
}

func UpdateCustomer(ctx context.Context, executor qrbpgx.Executor, id uuid.UUID, changeSet CustomerChangeSet) error {
	q := qrb.
		Update(customer).
		Where(customer.ID.Eq(qrb.Arg(id))).
		SetMap(changeSet.toMap())

	res, err := qrbpgx.Build(q).WithExecutor(executor).Exec(ctx)
	if err != nil {
		return err
	}

	return assertRowsAffected(res, "update", 1)
}

func customerJson() builder.JsonBuildObjectBuilder {
	// customerDefaultJson is generated by construct and is a JsonBuildObjectBuilder that can be further modified (immutable)
	return customerDefaultJson.
		// It's easy to set additional properties
		Prop("ProjectCount", qrb.Count(project.ID))
}
```

These are common functions that can be shared by all repository implementations:

*repository/commong.go*
```go
package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/networkteam/construct/v2"
)

func pgxCollectRow[T any](row pgx.CollectableRow) (T, error) {
	return pgxScanRow[T](row)
}

func pgxScanRow[T any](row pgx.Row) (T, error) {
	var result T
	err := row.Scan(&result)
	if err != nil {
		return result, fmt.Errorf("scanning row: %w", pgxToConstructErr(err))
	}
	return result, nil
}

func pgxToConstructErr(err error) error {
	if err == pgx.ErrNoRows {
		return construct.ErrNotFound
	}
	return err
}

func assertRowsAffected(res pgconn.CommandTag, op string, numberOfRows int64) error {
	rowsAffected := res.RowsAffected()
	if rowsAffected != numberOfRows {
		return fmt.Errorf("%s affected %d rows, but expected exactly %d", op, rowsAffected, numberOfRows)
	}
	return nil
}
```

## Install

Add Go module:

```bash
go get github.com/networkteam/construct/v2
```

Add module to a `tools.go` file:

```go
// +build tools
package myproject

import (
    _ "github.com/networkteam/construct/v2"
)
```

## License

[MIT](./LICENSE)
