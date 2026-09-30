package camelstring

import (
	"strings"
	"unicode"
)

// CamelToKebab converts camelCase or PascalCase to kebab-case.
// It strictly converts any uppercase letter to a hyphen followed by the lowercase letter.
func CamelToKebab(s string) string {
	var builder strings.Builder
	// Pre-allocate memory to reduce allocations
	builder.Grow(len(s) + 5)

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				builder.WriteRune('-')
			}
			builder.WriteRune(unicode.ToLower(r))
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func ToCamelUpper(s string) string {
	ret := ToCamel(s)
	if ret == "" {
		return ""
	}
	return strings.ToUpper(ret[0:1]) + ret[1:]
}

// ToCamel converts to lowerCamelCase.
func ToCamel(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))

	nextUpper := false
	for i, r := range s {
		if r == '-' || r == '_' {
			nextUpper = true
			continue
		}

		if nextUpper {
			builder.WriteRune(unicode.ToUpper(r))
			nextUpper = false
		} else if i == 0 {
			// Ensure the first character is lowercase for standard camelCase
			builder.WriteRune(unicode.ToLower(r))
		} else {
			// Normalize remaining characters
			builder.WriteRune(unicode.ToLower(r))
		}
	}
	return builder.String()
}
