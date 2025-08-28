package stringutil

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID creates a new UUID v7 and returns it as a string.
func GenerateUUID() (string, error) {
	value, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("uuid generation failed: %w", err)
	}

	return value.String(), nil
}

// ToSnakeCase converts a given string from CamelCase or PascalCase to snake_case.
func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ReplaceAll(snake, "__", "_")

	return strings.ToLower(snake)
}
