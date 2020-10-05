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
		// current index
		i int
		// index of rune before
		j, k     int
		hasLower bool
		// current rune
		r rune
	)
	for i, r = range s {
		k = j
		j = i
		if unicode.IsLower(r) {
			hasLower = true
			// Get the index before the last consecutive uppercase prefix rune (huh???)
			if k > 0 {
				j = k
			}
			break
		}
	}
	if !hasLower {
		j = len(s)
	}

	return strings.ToLower(s[0:j]) + s[j:]
}
