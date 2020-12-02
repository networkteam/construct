package json

import (
	"sort"
	"strings"

	"github.com/rsms/go-immutable"
)

// JsonBuildObjectBuilder is a builder for a JSON_BUILD_OBJECT function. It is immutable and safe for concurrent use.
type JsonBuildObjectBuilder struct {
	m *immutable.StrMap
}

var _ SqlGenerator = JsonBuildObjectBuilder{}

// JsonBuildObject starts a builder for a JSON_BUILD_OBJECT function
func JsonBuildObject() JsonBuildObjectBuilder {
	return JsonBuildObjectBuilder{m: immutable.EmptyStrMap}
}

// GenerateSql implements SqlGenerator
func (b JsonBuildObjectBuilder) GenerateSql(sb *strings.Builder) {
	sb.WriteString("JSON_BUILD_OBJECT(")
	i := 0

	keys := make([]string, 0, b.m.Len)
	b.m.Range(func(key string, v interface{}) bool {
		keys = append(keys, key)
		return true
	})
	sort.Strings(keys)

	for _, key := range keys {
		v := b.m.Get(key).(SqlGenerator)
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(pqQuoteString(key))
		sb.WriteRune(',')
		v.GenerateSql(sb)
		i++
	}
	sb.WriteString(")")
}

// Set returns a new builder with added or updated mapping from a property name to a value (SQL generator like Exp).
func (b JsonBuildObjectBuilder) Set(propertyName string, value SqlGenerator) JsonBuildObjectBuilder {
	return JsonBuildObjectBuilder{
		m: b.m.Set(propertyName, value),
	}
}

// Delete returns a new builder with removed property name from the mapping
func (b JsonBuildObjectBuilder) Delete(propertyName string) JsonBuildObjectBuilder {
	return JsonBuildObjectBuilder{
		m: b.m.Del(propertyName),
	}
}

// ToSql returns SQL for the builder
func (b JsonBuildObjectBuilder) ToSql() string {
	return GeneratorToSql(b)
}

func pqQuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
