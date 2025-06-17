package apperror

import (
	"errors"
)

type ErrorCode string

const (
	INVALID_ACTION_CODE        ErrorCode = "INVALID_ACTION"
	INVALID_DATA_CODE          ErrorCode = "INVALID_DATA"
	INVALID_PARAMS_CODE        ErrorCode = "INVALID_PARAMS"
	UNAUTHORIZED_CODE          ErrorCode = "UNAUTHORIZED"
	FORBIDDEN_CODE             ErrorCode = "FORBIDDEN"
	CSRF_TOKEN_MISMATCH_CODE   ErrorCode = "CSRF_MISMATCH"
	RECORD_NOT_FOUND_CODE      ErrorCode = "RECORD_NOT_FOUND"
	INTERNAL_SERVER_ERROR_CODE ErrorCode = "INTERNAL_SERVER_ERROR"
)

var (
	Err400InvalidAction     = errors.New("invalid action")
	Err400InvalidData       = errors.New("invalid data")
	Err400InvalidParams     = errors.New("invalid params")
	Err401Unauthorized      = errors.New("unauthorized")
	Err403Forbidden         = errors.New("forbidden")
	Err403CSRFTokenMismatch = errors.New("csrf token mismatch")
	Err404RecordNotFound    = errors.New("record not found")
	Err500InternalServer    = errors.New("internal server error")
)
