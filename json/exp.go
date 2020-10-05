package json

import (
	"strings"
)

type expGenerator string

func (e expGenerator) GenerateSql(sb *strings.Builder) {
	sb.WriteString(string(e))
}

// Exp generates SQL for an arbitrary expression
func Exp(exp string) SqlGenerator {
	return expGenerator(exp)
}
