package json

import (
	"strings"
)

type expGenerator string

// GenerateSql implements SqlGenerator
func (e expGenerator) GenerateSql(sb *strings.Builder) {
	sb.WriteString(string(e))
}

// Exp generates SQL for an arbitrary SQL expression
func Exp(exp string) SqlGenerator {
	return expGenerator(exp)
}
