package templit

import (
	"strings"
	"text/template"
)

// DefaultFuncMap is the default function map for templates.
var DefaultFuncMap = template.FuncMap{
	"lower":        strings.ToLower,
	"upper":        strings.ToUpper,
	"trim":         strings.TrimSpace,
	"split":        strings.Split,
	"join":         strings.Join,
	"replace":      strings.ReplaceAll,
	"contains":     strings.Contains,
	"hasPrefix":    strings.HasPrefix,
	"hasSuffix":    strings.HasSuffix,
	"trimPrefix":   strings.TrimPrefix,
	"trimSuffix":   strings.TrimSuffix,
	"trimSpace":    strings.TrimSpace,
	"trimLeft":     strings.TrimLeft,
	"trimRight":    strings.TrimRight,
	"count":        strings.Count,
	"repeat":       strings.Repeat,
	"equalFold":    strings.EqualFold,
	"splitN":       strings.SplitN,
	"splitAfter":   strings.SplitAfter,
	"splitAfterN":  strings.SplitAfterN,
	"fields":       strings.Fields,
	"toTitle":      strings.ToTitle,
	"toSnakeCase":  ToSnakeCase,
	"toCamelCase":  ToCamelCase,
	"toKebabCase":  ToKebabCase,
	"toPascalCase": ToPascalCase,
	"default":      defaultVal,
}

// defaultVal returns defaultValue if value is nil, otherwise value.
func defaultVal(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	return value
}
