package structutil

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shoraid/stx-go-utils/apperror"
	"github.com/stretchr/testify/assert"
)

func TestStructUtil_Validate(t *testing.T) {
	type UserRequest struct {
		Name     string `json:"name" validate:"required,max=10"`
		Email    string `json:"email" validate:"required,email"`
		Age      int    `json:"age" validate:"min=18"`
		IsActive bool   `json:"is_active" validate:"boolean"`
	}

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
				"email": {"field must be a valid email address"},
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
			result, err := Validate(tt.request)

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

func TestStructUtil_Validate_Nested(t *testing.T) {
	type Meta struct {
		Note string `json:"note" validate:"required"`
	}

	type Roles struct {
		ID   string `json:"id" validate:"required,uuid"`
		Name string `json:"name" validate:"required"`
	}

	type UserRequest struct {
		Meta          Meta     `json:"meta" validate:"required"`
		Roles         []Roles  `json:"roles" validate:"required,dive"`
		PermissionIDs []string `json:"permissionIds" validate:"required,dive,uuid"`
	}

	tests := []struct {
		name     string
		request  any
		expected map[string][]string
		isError  bool
	}{
		{
			name: "Valid nested request",
			request: UserRequest{
				Meta: Meta{Note: "Valid note"},
				Roles: []Roles{
					{ID: "d290f1ee-6c54-4b01-90e6-d701748f0851", Name: "Admin"},
				},
				PermissionIDs: []string{"d290f1ee-6c54-4b01-90e6-d701748f0851"},
			},
			expected: nil,
			isError:  false,
		},
		{
			name: "Empty meta.note",
			request: UserRequest{
				Meta:          Meta{Note: ""},
				Roles:         []Roles{{ID: "d290f1ee-6c54-4b01-90e6-d701748f0851", Name: "Admin"}},
				PermissionIDs: []string{"d290f1ee-6c54-4b01-90e6-d701748f0851"},
			},
			expected: map[string][]string{
				"meta.note": {"field is required"},
			},
			isError: true,
		},
		{
			name: "Invalid role ID and empty role name",
			request: UserRequest{
				Roles: []Roles{
					{ID: "invalid-uuid", Name: ""},
				},
				Meta:          Meta{Note: "Valid note"},
				PermissionIDs: []string{"d290f1ee-6c54-4b01-90e6-d701748f0851"},
			},
			expected: map[string][]string{
				"roles.0.id":   {"field must be a valid UUID"},
				"roles.0.name": {"field is required"},
			},
			isError: true,
		},
		{
			name: "Invalid permission ID",
			request: UserRequest{
				PermissionIDs: []string{"invalid-uuid"},
				Meta:          Meta{Note: "Valid note"},
				Roles:         []Roles{{ID: "d290f1ee-6c54-4b01-90e6-d701748f0851", Name: "Admin"}},
			},
			expected: map[string][]string{
				"permissionIds.0": {"field must be a valid UUID"},
			},
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Validate(tt.request)

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
			result, err := BindAndValidateJSON(req, &input)

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

func TestStructUtil_getErrorMessage(t *testing.T) {
	type Sample struct {
		Email    string `validate:"email"`
		Name     string `validate:"required"`
		Age      int    `validate:"min=18"`
		IsActive bool   `validate:"boolean"`
		Role     string `validate:"oneof=admin user"`
		ID       string `validate:"uuid"`
		MaxTest  string `validate:"max=5"`
	}

	validate := validator.New()
	s := Sample{
		Email:    "invalid-email",
		Name:     "",
		Age:      15,
		IsActive: true,
		Role:     "guest",
		ID:       "invalid-uuid",
		MaxTest:  "toolongstring",
	}

	err := validate.Struct(s)
	assert.Error(t, err)

	validationErrors := err.(validator.ValidationErrors)

	tests := []struct {
		tag      string
		expected string
	}{
		{"required", "field is required"},
		{"email", "field must be a valid email address"},
		{"min", "minimum value is 18"},
		{"boolean", "field must be a boolean"}, // no actual error here, but still tested
		{"oneof", "field must be one of: admin, user"},
		{"uuid", "field must be a valid UUID"},
		{"max", "maximum length is 5"},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			// find the matching error with the tag
			var matched validator.FieldError
			for _, fe := range validationErrors {
				if fe.Tag() == tt.tag {
					matched = fe
					break
				}
			}
			// only test if matched tag exists in this struct
			if matched != nil {
				actual := getErrorMessage(matched)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestStructUtil_getJSONTagName(t *testing.T) {
	type TestStruct struct {
		WithTag      string `json:"with_tag"`
		WithTagOmit  string `json:"with_tag_omit,omitempty"`
		WithoutTag   string
		IgnoredField string `json:"-"`
		EmptyTag     string `json:""`
	}

	tests := []struct {
		fieldName string
		expected  string
	}{
		{"WithTag", "with_tag"},
		{"WithTagOmit", "with_tag_omit"},
		{"WithoutTag", "WithoutTag"},
		{"IgnoredField", "IgnoredField"},
		{"EmptyTag", "EmptyTag"},
	}

	tType := reflect.TypeOf(TestStruct{})

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			field, ok := tType.FieldByName(tt.fieldName)
			assert.True(t, ok, "field should exist")

			result := getJSONTagName(field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func BenchmarkStructutil_Validate(b *testing.B) {
	type UserRequest struct {
		Name     string `json:"name" validate:"required,max=10"`
		Email    string `json:"email" validate:"required,email"`
		Age      int    `json:"age" validate:"min=18"`
		IsActive bool   `json:"is_active" validate:"boolean"`
	}

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
			for b.Loop() {
				Validate(tt.payload)
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
			for b.Loop() {
				req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")

				var input LoginRequest
				BindAndValidateJSON(req, &input)
			}
		})
	}
}

func BenchmarkStructutil_getErrorMessage(b *testing.B) {
	type BenchmarkSample struct {
		Email    string `validate:"email"`
		Name     string `validate:"required"`
		Age      int    `validate:"min=18"`
		IsActive string `validate:"boolean"`
		Role     string `validate:"oneof=admin user"`
		ID       string `validate:"uuid"`
		MaxTest  string `validate:"max=5"`
	}

	validate := validator.New()
	s := BenchmarkSample{
		Email:    "not-an-email",
		Name:     "",
		Age:      10,
		IsActive: "maybe",
		Role:     "guest",
		ID:       "invalid-uuid",
		MaxTest:  "this string is too long",
	}

	err := validate.Struct(s)
	if err == nil {
		b.Fatal("Expected validation error, got nil")
	}

	errors := err.(validator.ValidationErrors)

	for b.Loop() {
		for _, fe := range errors {
			getErrorMessage(fe)
		}
	}
}

func BenchmarkStructutil_getJSONTagName(b *testing.B) {
	type BenchmarkStruct struct {
		WithTag      string `json:"with_tag"`
		WithTagOmit  string `json:"with_tag_omit,omitempty"`
		WithoutTag   string
		IgnoredField string `json:"-"`
		EmptyTag     string `json:""`
	}

	tType := reflect.TypeOf(BenchmarkStruct{})
	fields := []reflect.StructField{
		tType.Field(0), // WithTag
		tType.Field(1), // WithTagOmit
		tType.Field(2), // WithoutTag
		tType.Field(3), // IgnoredField
		tType.Field(4), // EmptyTag
	}

	for b.Loop() {
		for _, field := range fields {
			getJSONTagName(field)
		}
	}
}
