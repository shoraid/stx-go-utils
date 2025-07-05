package stringutil

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ReplaceAll(snake, "__", "_")

	return strings.ToLower(snake)
}

func GenerateUUID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("uuid generation failed: %w", err)
	}

	return id.String(), nil
}
