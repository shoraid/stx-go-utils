package httpresponse

import (
	"encoding/json"
	"net/http"

	"github.com/shoraid/stx-go-utils/apperror"
)

type Response struct {
	Code    apperror.ErrorCode `json:"code"`
	Message string             `json:"message"`
	Details any                `json:"details,omitempty"`
}

func HandleError(w http.ResponseWriter, err error, details ...any) bool {
	if err == nil {
		return false
	}

	var errorDetails any
	if len(details) > 0 {
		errorDetails = details[0]
	}

	var resp Response
	var statusCode int

	switch err {
	case apperror.Err400InvalidAction:
		resp = Response{
			Code:    apperror.INVALID_ACTION_CODE,
			Message: "Invalid action",
			Details: errorDetails,
		}
		statusCode = http.StatusBadRequest

	case apperror.Err400InvalidData:
		resp = Response{
			Code:    apperror.INVALID_DATA_CODE,
			Message: "Invalid data",
			Details: map[string]any{"errors": errorDetails},
		}
		statusCode = http.StatusBadRequest

	case apperror.Err400InvalidBody:
		resp = Response{
			Code:    apperror.INVALID_BODY_CODE,
			Message: "Invalid body",
			Details: map[string]any{"errors": errorDetails},
		}
		statusCode = http.StatusBadRequest

	case apperror.Err400InvalidParams:
		resp = Response{
			Code:    apperror.INVALID_PARAMS_CODE,
			Message: "Invalid params",
			Details: errorDetails,
		}
		statusCode = http.StatusBadRequest

	case apperror.Err401Unauthorized:
		resp = Response{
			Code:    apperror.UNAUTHORIZED_CODE,
			Message: "Unauthorized",
			Details: errorDetails,
		}
		statusCode = http.StatusUnauthorized

	case apperror.Err403Forbidden:
		resp = Response{
			Code:    apperror.FORBIDDEN_CODE,
			Message: "Forbidden",
			Details: errorDetails,
		}
		statusCode = http.StatusForbidden

	case apperror.Err403CSRFTokenMismatch:
		resp = Response{
			Code:    apperror.CSRF_TOKEN_MISMATCH_CODE,
			Message: "CSRF token mismatch",
			Details: errorDetails,
		}
		statusCode = http.StatusForbidden

	case apperror.Err404RecordNotFound:
		resp = Response{
			Code:    apperror.RECORD_NOT_FOUND_CODE,
			Message: "Record not found",
			Details: errorDetails,
		}
		statusCode = http.StatusNotFound

	default:
		resp = Response{
			Code:    apperror.INTERNAL_SERVER_ERROR_CODE,
			Message: "Internal server error",
			Details: errorDetails,
		}
		statusCode = http.StatusInternalServerError
	}

	writeJSON(w, statusCode, resp)
	return true
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
