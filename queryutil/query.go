package queryutil

import (
	"strconv"
	"strings"
)

func CalculatePagination(pageStr, perPageStr string, defaultPerPage int) (page, perPage, offset int) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	perPage, err = strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		perPage = defaultPerPage
	}

	offset = (page - 1) * perPage

	return page, perPage, offset
}

func CalculateTotalPage(totalData int, perPage int) int {
	// Validate both totalData and perPage
	if totalData < 0 || perPage <= 0 {
		return 0
	}

	totalPage := totalData / perPage

	if totalData%perPage > 0 {
		// Round up if there's a remainder
		totalPage++
	}

	return totalPage
}

// ResolveAllowedFields parses a comma-separated input string and returns only the
// items that are allowed based on the provided map.
//
// Parameters:
// - input: comma-separated string (e.g., "name,email").
// - allowed: map of allowed fields, where each key can be:
//   - bool (true): to allow the field as-is.
//   - string: to alias the field to a different value.
//
// Returns:
// - A slice of strings that are allowed according to the map.
//
// Example:
//
//	input = "name,email"
//	allowed = map[string]any{"name": true, "email": "user_email"}
//	→ returns: []string{"name", "user_email"}
func ResolveAllowedFields(input string, allowed map[string]any) []string {
	if input == "" {
		return []string{}
	}

	splitted := strings.Split(input, ",")
	result := make([]string, 0, len(splitted))

	for _, item := range splitted {
		field := strings.TrimSpace(item)

		if val, ok := allowed[field]; ok {
			switch v := val.(type) {
			case bool:
				if v {
					result = append(result, field)
				}
			case string:
				result = append(result, v)
			}
		}
	}

	return result
}

// ResolveSingleField checks if the input exists in the allowed map and returns the mapped value
// if available, otherwise returns the defaultField.
//
// Parameters:
// - input: the field name to check (e.g., "email").
// - allowed: map of allowed fields, where each key can be:
//   - bool (true): to allow the field as-is.
//   - string: to alias the field to a different value.
//
// - defaultField: fallback value to return if input is not allowed.
//
// Returns:
// - A string representing the resolved field name if allowed, or the defaultField if not.
//
// Example:
//
//	input = "email"
//	allowed = map[string]any{"email": "user_email", "username": true}
//	defaultField = "username"
//	→ returns: "user_email"
func ResolveSingleField(input string, allowed map[string]any, defaultField string) string {
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultField
	}

	if val, ok := allowed[input]; ok {
		switch v := val.(type) {
		case bool:
			if v {
				return input
			}
		case string:
			return v
		}
	}

	return defaultField
}
