module github.com/networkteam/construct/v2/example

go 1.15

require (
	github.com/Masterminds/squirrel v1.4.0
	github.com/friendsofgo/errors v0.9.2
	github.com/gofrs/uuid v4.3.0+incompatible
	github.com/jackc/pgconn v1.7.0
	github.com/jackc/pgtype v1.5.0
	github.com/jackc/pgx/v4 v4.9.0 // indirect
	github.com/networkteam/construct/v2 v2.0.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/urfave/cli/v2 v2.17.1
)

replace github.com/networkteam/construct/v2 => ../
