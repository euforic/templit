package templit

import (
	"strings"
	"text/template"
)

// DefaultFuncMap is the default function map for templates.
func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":         strings.ToLower,
		"upper":         strings.ToUpper,
		"trim":          strings.TrimSpace,
		"split":         strings.Split,
		"join":          strings.Join,
		"replace":       strings.ReplaceAll,
		"contains":      strings.Contains,
		"has_prefix":    strings.HasPrefix,
		"has_suffix":    strings.HasSuffix,
		"trim_prefix":   strings.TrimPrefix,
		"trim_suffix":   strings.TrimSuffix,
		"trim_space":    strings.TrimSpace,
		"trim_left":     strings.TrimLeft,
		"trim_right":    strings.TrimRight,
		"count":         strings.Count,
		"repeat":        strings.Repeat,
		"equal_fold":    strings.EqualFold,
		"split_n":       strings.SplitN,
		"split_after":   strings.SplitAfter,
		"split_after_n": strings.SplitAfterN,
		"fields":        strings.Fields,
		"title_case":    strings.ToTitle,
		"snake_case":    ToSnakeCase,
		"camel_case":    ToCamelCase,
		"kebab_case":    ToKebabCase,
		"pascal_case":   ToPascalCase,
		"default":       defaultVal,
	}
}

// defaultVal returns defaultValue if value is nil, otherwise value.
func defaultVal(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}

	return value
}
