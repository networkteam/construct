package json

import (
	"strings"
)

// JsonBuildArrayBuilder is a builder for a JSON_BUILD_ARRAY function. It is immutable and safe for concurrent use.
type JsonBuildArrayBuilder struct {
	m *sqlGenList
}

var _ SqlGenerator = JsonBuildArrayBuilder{}

// JsonBuildArray starts a builder for a JSON_BUILD_ARRAY function
func JsonBuildArray() JsonBuildArrayBuilder {
	return JsonBuildArrayBuilder{m: &sqlGenList{}}
}

// GenerateSql implements SqlGenerator
func (b JsonBuildArrayBuilder) GenerateSql(sb *strings.Builder) {
	sb.WriteString("JSON_BUILD_ARRAY(")
	i := 0

	b.m.iterate(func(v SqlGenerator) {
		if i > 0 {
			sb.WriteRune(',')
		}
		v.GenerateSql(sb)
		i++
	})
	sb.WriteString(")")
}

// Append returns a new builder with added array entry at the end for a value (SQL generator like Exp)
func (b JsonBuildArrayBuilder) Append(value SqlGenerator) JsonBuildArrayBuilder {
	return JsonBuildArrayBuilder{
		m: b.m.append(value),
	}
}

// ToSql returns SQL for the builder
func (b JsonBuildArrayBuilder) ToSql() string {
	return GeneratorToSql(b)
}
