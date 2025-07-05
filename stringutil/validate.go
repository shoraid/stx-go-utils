package stringutil

import "github.com/google/uuid"

func IsValidUUID(rawUUID string) bool {
	return uuid.Validate(rawUUID) == nil
}
