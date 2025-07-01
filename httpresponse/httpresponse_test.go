package httpresponse_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shoraid/stx-go-utils/apperror"
	"github.com/shoraid/stx-go-utils/httpresponse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpResponse_HandleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		details        []any
		expectedCode   int
		expectedBody   map[string]any
		expectedReturn bool
	}{
		{
			name:           "no error should return false",
			err:            nil,
			expectedCode:   http.StatusOK,
			expectedBody:   nil,
			expectedReturn: false,
		},
		{
			name:         "invalid action should return 400",
			err:          apperror.Err400InvalidAction,
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]any{
				"code":    string(apperror.INVALID_ACTION_CODE),
				"message": "Invalid action",
				"details": nil,
			},
			expectedReturn: true,
		},
		{
			name:         "invalid body should return 400",
			err:          apperror.Err400InvalidBody,
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]any{
				"code":    string(apperror.INVALID_BODY_CODE),
				"message": "Invalid body",
				"details": map[string]any{
					"errors": nil,
				},
			},
			expectedReturn: true,
		},
		{
			name:         "invalid data should return 400",
			err:          apperror.Err400InvalidData,
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]any{
				"code":    string(apperror.INVALID_DATA_CODE),
				"message": "Invalid data",
				"details": map[string]any{
					"errors": nil,
				},
			},
			expectedReturn: true,
		},
		{
			name:         "invalid params should return 400",
			err:          apperror.Err400InvalidParams,
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]any{
				"code":    string(apperror.INVALID_PARAMS_CODE),
				"message": "Invalid params",
				"details": nil,
			},
			expectedReturn: true,
		},
		{
			name:         "unauthorized should return 401",
			err:          apperror.Err401Unauthorized,
			expectedCode: http.StatusUnauthorized,
			expectedBody: map[string]any{
				"code":    string(apperror.UNAUTHORIZED_CODE),
				"message": "Unauthorized",
				"details": nil,
			},
			expectedReturn: true,
		},
		{
			name:         "forbidden should return 403",
			err:          apperror.Err403Forbidden,
			expectedCode: http.StatusForbidden,
			expectedBody: map[string]any{
				"code":    string(apperror.FORBIDDEN_CODE),
				"message": "Forbidden",
				"details": nil,
			},
			expectedReturn: true,
		},
		{
			name:         "csrf token mismatch should return 403",
			err:          apperror.Err403CSRFTokenMismatch,
			expectedCode: http.StatusForbidden,
			expectedBody: map[string]any{
				"code":    string(apperror.CSRF_TOKEN_MISMATCH_CODE),
				"message": "CSRF token mismatch",
				"details": nil,
			},
			expectedReturn: true,
		},
		{
			name:         "record not found should return 404",
			err:          apperror.Err404RecordNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]any{
				"code":    string(apperror.RECORD_NOT_FOUND_CODE),
				"message": "Record not found",
				"details": nil,
			},
			expectedReturn: true,
		},
		{
			name:         "internal server error should return 500",
			err:          apperror.Err500InternalServer,
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"code":    string(apperror.INTERNAL_SERVER_ERROR_CODE),
				"message": "Internal server error",
				"details": nil,
			},
			expectedReturn: true,
		},
		{
			name:         "default error should return 500",
			err:          errors.New("default error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]any{
				"code":    string(apperror.INTERNAL_SERVER_ERROR_CODE),
				"message": "Internal server error",
				"details": nil,
			},
			expectedReturn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			got := httpresponse.HandleError(rec, tt.err, tt.details...)

			assert.Equal(t, tt.expectedReturn, got)

			if !tt.expectedReturn {
				assert.Equal(t, tt.expectedCode, rec.Code)
				return
			}

			require.Equal(t, tt.expectedCode, rec.Code)

			var resp map[string]any
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			require.NoError(t, err)

			for key, expected := range tt.expectedBody {
				assert.Equal(t, expected, resp[key], "Mismatch on key: %s", key)
			}
		})
	}
}

func BenchmarkHttpResponse_HandleError(b *testing.B) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "NoError",
			err:  nil,
		},
		{
			name: "InvalidActionError",
			err:  apperror.Err400InvalidAction,
		},
		{
			name: "InvalidDataError",
			err:  apperror.Err400InvalidData,
		},
		{
			name: "InvalidBodyError",
			err:  apperror.Err400InvalidBody,
		},
		{
			name: "InvalidParamsError",
			err:  apperror.Err400InvalidParams,
		},
		{
			name: "UnauthorizedError",
			err:  apperror.Err401Unauthorized,
		},
		{
			name: "ForbiddenError",
			err:  apperror.Err403Forbidden,
		},
		{
			name: "CSRFTokenMismatchError",
			err:  apperror.Err403CSRFTokenMismatch,
		},
		{
			name: "AppRecordNotFoundError",
			err:  apperror.Err404RecordNotFound,
		},
		{
			name: "InternalServerError",
			err:  apperror.Err500InternalServer,
		},
		{
			name: "DefaultError",
			err:  errors.New("default error"),
		},
	}

	b.ResetTimer()
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				_ = httpresponse.HandleError(w, tt.err)
			}
		})
	}
}
