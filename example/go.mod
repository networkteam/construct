module github.com/networkteam/construct/example

go 1.15

require (
	github.com/Masterminds/squirrel v1.4.0
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/friendsofgo/errors v0.9.2
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/jackc/pgconn v1.7.0
	github.com/jackc/pgtype v1.5.0
	github.com/jackc/pgx/v4 v4.9.0
	github.com/networkteam/construct v0.0.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/urfave/cli/v2 v2.2.0
)

replace github.com/networkteam/construct => ../
