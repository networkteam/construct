package json

import "strings"

// SqlGenerator is an interface for all types that generate SQL by appending to a given strings.Builder
type SqlGenerator interface {
	GenerateSql(sb *strings.Builder)
}
