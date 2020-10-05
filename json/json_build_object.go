package json

import "strings"

type SqlGenerator interface {
	GenerateSql(sb *strings.Builder)
}

type JsonBuildObjectBuilder map[string]SqlGenerator

var _ SqlGenerator = JsonBuildObjectBuilder{}

func (j JsonBuildObjectBuilder) GenerateSql(sb *strings.Builder) {
	sb.WriteString("JSON_BUILD_OBJECT(")
	i := 0

	keys := make([]string, 0, len(j))
	for key := range j {
		keys = append(keys, key)
	}

	for _, key := range keys {
		v := j[key]
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

// JsonBuildObject generates SQL for a JSON_BUILD_OBJECT function
func JsonBuildObject() JsonBuildObjectBuilder {
	return make(JsonBuildObjectBuilder)
}

func (b JsonBuildObjectBuilder) Set(key string, value SqlGenerator) JsonBuildObjectBuilder {
	b[key] = value
	return b
}

func (j JsonBuildObjectBuilder) Delete(key string) {
	delete(j, key)
}

func (j JsonBuildObjectBuilder) ToSql() string {
	sb := new(strings.Builder)
	j.GenerateSql(sb)
	return sb.String()
}

func pqQuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
