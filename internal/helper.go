package internal

import (
	"strings"
	"unicode"
)

func firstToUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func firstToLower(s string) string {
	if s == "" {
		return ""
	}
	var (
		i int
		secondIdx int
		r rune
	)
	// Check if all runes are uppercase
	for i, r = range s {
		if i > 0 && secondIdx == 0 {
			secondIdx = i
		}
		if unicode.IsLower(r) {
			break
		}
	}
	// All runes are uppercase (e.g. "ID")
	if i == len(s)-1 {
		return strings.ToLower(s)
	}
	// Only lowercase first rune
	return strings.ToLower(s[0:secondIdx]) + s[secondIdx:]
}
