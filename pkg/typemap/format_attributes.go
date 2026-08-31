package typemap

import (
	"strconv"
	"strings"
)

type FormatHints struct {
	Format    string
	MinLength int
	MaxLength int
	Regex     string
	Min       float64
	Max       float64
	HasMin    bool
	HasMax    bool
}

func ParseFormatHints(attributes []string) FormatHints {
	hints := FormatHints{}
	for _, attr := range attributes {
		key, value, found := strings.Cut(attr, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "format":
			hints.Format = value
		case "minLength":
			if v, err := strconv.Atoi(value); err == nil {
				hints.MinLength = v
			}
		case "maxLength":
			if v, err := strconv.Atoi(value); err == nil {
				hints.MaxLength = v
			}
		case "regex":
			hints.Regex = value
		case "min":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				hints.Min = v
				hints.HasMin = true
			}
		case "max":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				hints.Max = v
				hints.HasMax = true
			}
		}
	}
	return hints
}

func HasAttribute(attributes []string, target string) bool {
	for _, attr := range attributes {
		if attr == target {
			return true
		}
	}
	return false
}
