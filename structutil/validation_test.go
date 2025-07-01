package structutil_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shoraid/stx-go-utils/apperror"
	"github.com/shoraid/stx-go-utils/structutil"
	"github.com/stretchr/testify/assert"
)

type UserRequest struct {
	Name     string `json:"name" validate:"required,max=10"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"min=18"`
	IsActive bool   `json:"is_active" validate:"boolean"`
}

func TestStructUtil_Validate(t *testing.T) {
	tests := []struct {
		name     string
		request  any
		expected map[string][]string
		isError  bool
	}{
		{
			name: "Valid struct should return no error",
			request: UserRequest{
				Name:     "Alice",
				Email:    "alice@example.com",
				Age:      25,
				IsActive: true,
			},
			expected: nil,
			isError:  false,
		},
		{
			name: "Missing required fields",
			request: UserRequest{
				Name:     "",
				Email:    "",
				Age:      25,
				IsActive: true,
			},
			expected: map[string][]string{
				"name":  {"field is required"},
				"email": {"field is required"},
			},
			isError: true,
		},
		{
			name: "Exceeds max name length",
			request: UserRequest{
				Name:     "ThisNameIsWayTooLong",
				Email:    "john@example.com",
				Age:      25,
				IsActive: true,
			},
			expected: map[string][]string{
				"name": {"maximum length is 10"},
			},
			isError: true,
		},
		{
			name: "Age below minimum",
			request: UserRequest{
				Name:     "John",
				Email:    "john@example.com",
				Age:      17,
				IsActive: true,
			},
			expected: map[string][]string{
				"age": {"minimum value is 18"},
			},
			isError: true,
		},
		{
			name: "Multiple validation errors",
			request: UserRequest{
				Name:     "",
				Email:    "invalid-email",
				Age:      15,
				IsActive: true,
			},
			expected: map[string][]string{
				"name":  {"field is required"},
				"email": {"field is invalid"},
				"age":   {"minimum value is 18"},
			},
			isError: true,
		},
		{
			name: "Pointer struct with validation errors",
			request: &UserRequest{
				Name:     "",
				Email:    "",
				Age:      10,
				IsActive: true,
			},
			expected: map[string][]string{
				"name":  {"field is required"},
				"email": {"field is required"},
				"age":   {"minimum value is 18"},
			},
			isError: true,
		},
		{
			name: "Struct with nil pointer fields",
			request: struct {
				Description *string `json:"description" validate:"required"`
				Count       *int    `json:"count" validate:"required,min=1"`
				Active      *bool   `json:"active" validate:"required,boolean"`
			}{
				Description: nil,
				Count:       nil,
				Active:      nil,
			},
			expected: map[string][]string{
				"description": {"field is required"},
				"count":       {"field is required"},
				"active":      {"field is required"},
			},
			isError: true,
		},
		{
			name: "Field with no json tag should fallback to field name",
			request: struct {
				NoTag string `validate:"required"`
			}{
				NoTag: "",
			},
			expected: map[string][]string{
				"NoTag": {"field is required"},
			},
			isError: true,
		},
		{
			name: "Field with json:\"-\" should fallback to field name",
			request: struct {
				Ignored string `json:"-" validate:"required"`
			}{
				Ignored: "",
			},
			expected: map[string][]string{
				"Ignored": {"field is required"},
			},
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := structutil.Validate(tt.request)

			if tt.isError {
				assert.Equal(t, apperror.Err400InvalidData, err)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, result)
			}
		})
	}
}

func TestStructUtil_BindAndValidateJSON(t *testing.T) {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	tests := []struct {
		name           string
		body           string
		expectedError  error
		expectedFields map[string][]string
	}{
		{
			name: "Valid JSON and valid data",
			body: `{"email":"test@example.com","password":"secret123"}`,
		},
		{
			name:          "Valid JSON but invalid data",
			body:          `{"email":"","password":"123"}`,
			expectedError: apperror.Err400InvalidData,
			expectedFields: map[string][]string{
				"email":    {"field is required"},
				"password": {"minimum value is 6"},
			},
		},
		{
			name:          "Empty body (EOF)",
			body:          ``,
			expectedError: apperror.Err400InvalidData,
			expectedFields: map[string][]string{
				"email":    {"field is required"},
				"password": {"field is required"},
			},
		},
		{
			name:          "Invalid JSON (trailing comma)",
			body:          `{"email":"test@example.com","password":"123",}`,
			expectedError: apperror.Err400InvalidBody,
			expectedFields: map[string][]string{
				"json": {"invalid JSON format: please check for missing commas, braces, or quotes"},
			},
		},
		{
			name:          "Wrong type (password should be string)",
			body:          `{"email":"test@example.com","password":123}`,
			expectedError: apperror.Err400InvalidData,
			expectedFields: map[string][]string{
				"password": {"invalid type, expected string"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			var input LoginRequest
			result, err := structutil.BindAndValidateJSON(req, &input)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Equal(t, tt.expectedFields, result)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "test@example.com", input.Email)
				assert.Equal(t, "secret123", input.Password)
			}
		})
	}
}

func BenchmarkStructutil_Validate(b *testing.B) {
	valid := UserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Age:      25,
		IsActive: true,
	}

	invalid := UserRequest{
		Name:     "",             // required
		Email:    "invalidemail", // invalid
		Age:      10,             // < 18
		IsActive: true,
	}

	tests := []struct {
		name    string
		payload any
	}{
		{"ValidStruct", valid},
		{"InvalidStruct", invalid},
		{"PointerValidStruct", &valid},
		{"PointerInvalidStruct", &invalid},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				structutil.Validate(tt.payload)
			}
		})
	}
}

func BenchmarkStructutil_BindAndValidateJSON(b *testing.B) {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	tests := []struct {
		name string
		body string
	}{
		{
			"Valid JSON and valid data",
			`{"email":"test@example.com","password":"secret123"}`,
		},
		{
			"Valid JSON but invalid data",
			`{"email":"","password":"123"}`,
		},
		{
			"Empty body (EOF)",
			``,
		},
		{
			"Invalid JSON (trailing comma)",
			`{"email":"test@example.com","password":"123",}`,
		},
		{
			"Wrong type (password should be string)",
			`{"email":"test@example.com","password":123}`,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")

				var input LoginRequest
				_, _ = structutil.BindAndValidateJSON(req, &input)
			}
		})
	}
}
