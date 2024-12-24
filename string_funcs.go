package templit

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a string to CamelCase.
func ToCamelCase(s string) string {
	parts := splitAndFilter(s)
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = capitalizeFirstLetter(part)
		}
	}

	return strings.Join(parts, "")
}

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(s string) string {
	var result strings.Builder
	previousIsLower := false

	for _, r := range s {
		switch {
		case unicode.IsUpper(r) && previousIsLower:
			result.WriteRune('_')
			result.WriteRune(unicode.ToLower(r))
			previousIsLower = false
		case r == ' ', r == '-', r == '_':
			result.WriteRune('_')
			previousIsLower = false
		case unicode.IsLower(r) || unicode.IsDigit(r):
			result.WriteRune(r)
			previousIsLower = true
		default:
			result.WriteRune(unicode.ToLower(r))
			previousIsLower = false
		}
	}

	return result.String()
}

// ToKebabCase converts a string to kebab-case.
func ToKebabCase(s string) string {
	var result strings.Builder
	previousIsLower := false

	for _, r := range s {
		switch {
		case unicode.IsUpper(r) && previousIsLower:
			result.WriteRune('-')
			result.WriteRune(unicode.ToLower(r))
			previousIsLower = false
		case r == ' ', r == '_', r == '-':
			result.WriteRune('-')
			previousIsLower = false
		case unicode.IsLower(r) || unicode.IsDigit(r):
			result.WriteRune(r)
			previousIsLower = true
		default:
			result.WriteRune(unicode.ToLower(r))
			previousIsLower = false
		}
	}

	return result.String()
}

// ToPascalCase converts a string to PascalCase.
func ToPascalCase(s string) string {
	parts := splitAndFilter(s)
	for i, part := range parts {
		parts[i] = capitalizeFirstLetter(part)
	}

	return strings.Join(parts, "")
}

// splitAndFilter splits a string by multiple delimiters and filters out non-alphanumeric characters.
func splitAndFilter(s string) []string {
	parts := splitByMultipleDelimiters(s, []string{" ", "_", "-"})
	for i, part := range parts {
		runes := []rune(part)
		filtered := make([]rune, 0, len(runes))
		for _, r := range runes {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				filtered = append(filtered, r)
			}
		}
		parts[i] = string(filtered)
	}

	return parts
}

// splitByMultipleDelimiters splits a string based on multiple delimiters.
func splitByMultipleDelimiters(s string, delimiters []string) []string {
	for _, delimiter := range delimiters {
		s = strings.ReplaceAll(s, delimiter, " ")
	}

	return strings.Fields(s)
}

// capitalizeFirstLetter capitalizes the first letter of a string.
func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
