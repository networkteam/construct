# construct : Go generators for low abstraction persistence with PostgreSQL

[![GoDoc](https://godoc.org/github.com/networkteam/construct?status.svg)](https://godoc.org/github.com/networkteam/construct)
[![Build Status](https://github.com/networkteam/construct/workflows/run%20tests/badge.svg)](https://github.com/networkteam/construct/actions?workflow=run%20tests)
[![Coverage Status](https://coveralls.io/repos/github/networkteam/construct/badge.svg?branch=main)](https://coveralls.io/github/networkteam/construct?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/networkteam/construct)](https://goreportcard.com/report/github.com/networkteam/construct)

## Overview

Got tired of too many abstractions over all the features PostgreSQL provides when using an ORM? But rolling your own persistence code is tedious and there's too much boilderplate?
This is a code generator to generate a bunch of structs and functions to implement persistence code with a few line and keep all the power PostgreSQL provides.

## Example

**Your custom type defines fields for reading and writing**

*model/customer.go*
```go
package model

// Customer is just a struct, only struct tags are used in construct to generate code (no struct embedding needed) 
type Customer struct {
	ID            uuid.UUID `read_col:"customers.customer_id" write_col:"customer_id"`
	Name          string    `read_col:"customers.name,sortable" write_col:"name"`
    // Fields can be serialized as JSON by adding a "json" option to the "write_col" tag.
    // It works perfectly with a column of type jsonb or json in PostgreSQL.
	ContactPerson Contact   `read_col:"customers.contact_person,sortable" write_col:"contact_person,json"`

	// DomainCount is not mapped to the table but used in the select for reading an aggregate count
	DomainCount int

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

//go:generate go run github.com/networkteam/construct/cmd/construct my/project/model.Customer
```

```bash
go generate ./repository
```

**Roll your own persistence code for full control and low abstraction**

Structure your persistence code as it fits the project. A very simple and working approach is to have a bunch of
functions that operator on the target type for finding, inserting, updating and deleting. Add more complex queries
as you like in the same way. 

There's no need to use Squirrel, but it helps to build correct SQL queries instead of relying on string manipulation
and manually handling placeholders and arguments.

*repository/customer_repository.go*
```go
package repository

func FindCustomerByID(ctx context.Context, runner squirrel.BaseRunner, id uuid.UUID) (domain.Customer, error) {
	row := queryBuilder(runner).
		Select(buildCustomerJson()).
		From("customers").
		LeftJoin("domains ON (domains.customer_id = customers.customer_id)").
		// constants for read columns are generated by construct to prevent typos
		Where(squirrel.Eq{customer_id: id}).
		GroupBy(customer_id).
		QueryRowContext(ctx)
	// customerScanJsonRow is generated by construct and will scan a JSON row result into the target type 
	return customerScanJsonRow(row)
}

// CustomerChangeSet is generated by construct for handling partially filled models
func InsertCustomer(ctx context.Context, runner squirrel.BaseRunner, changeSet CustomerChangeSet) error {
	_, err := queryBuilder(runner).
		Insert("customers").
		// toMap is generated by construct 
		SetMap(changeSet.toMap()).
		ExecContext(ctx)
	return err
}

func UpdateCustomer(ctx context.Context, runner squirrel.BaseRunner, id uuid.UUID, changeSet CustomerChangeSet) error {
	res, err := queryBuilder(runner).
		Update("customers").
		Where(squirrel.Eq{customer_id: id}).
		SetMap(changeSet.toMap()).
		ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "executing update")
	}
	return assertRowsAffected(res, "update", 1)
}

func buildCustomerJson() string {
	// customerDefaultSelectJson is generated by construct and will generate a JSON_BUILD_OBJECT SQL expression
	// for returning a result that can be directly unmarshalled
	return customerDefaultSelectJson.
		// It's easy to set additional properties
		Set("DomainCount", cjson.Exp("COUNT(domains.domain_id)")).
		ToSql()
}
```

These are common functions that can be shared by all repository implementations:

*repository/commong.go*
```go
package repository

func queryBuilder(runner squirrel.BaseRunner) squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(runner)
}

func assertRowsAffected(res sql.Result, op string, nunberOfRows int64) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "getting affected rows")
	}
	if rowsAffected != nunberOfRows {
		return errors.Errorf("%s affected %d rows, but expected exactly %d", op, rowsAffected, nunberOfRows)
	}
	return err
}
```

## Install

```
go get github.com/networkteam/construct
```

## License

MIT.

