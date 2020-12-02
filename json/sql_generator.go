package json

import "strings"

// SqlGenerator is an interface for all types that generate SQL by appending to a given strings.Builder
type SqlGenerator interface {
	GenerateSql(sb *strings.Builder)
}

// GeneratorToSql generates SQL from a generator
func GeneratorToSql(g SqlGenerator) string {
	sb := new(strings.Builder)
	g.GenerateSql(sb)
	return sb.String()
}
