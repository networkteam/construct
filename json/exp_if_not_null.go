package json

import (
	"strings"
)

type expIfNotNullGenerator struct {
	conditionExp string
	then         SqlGenerator
}

// GenerateSql implements SqlGenerator
func (e expIfNotNullGenerator) GenerateSql(sb *strings.Builder) {
	sb.WriteString("CASE WHEN ")
	sb.WriteString(e.conditionExp)
	sb.WriteString(" IS NOT NULL THEN ")
	e.then.GenerateSql(sb)
	sb.WriteString(" END")
}

// ExpIfNotNull generates SQL for a conditional generator. It is particularly useful when nesting JSON_BUILD_OBJECT for
// eagerly fetched relations that can be NULL.
//
// It results in CASE WHEN [conditionExp] IS NOT NULL THEN [then] END. It will have a NULL result if conditionExp
// is NULL.
//
// Example:
//   JsonBuildObject().
//     Set("related", ExpIfNotNull("related_id", JsonBuildObject().Set("ID", Exp("related_id"))))
func ExpIfNotNull(conditionExp string, then SqlGenerator) SqlGenerator {
	return expIfNotNullGenerator{
		conditionExp,
		then,
	}
}
