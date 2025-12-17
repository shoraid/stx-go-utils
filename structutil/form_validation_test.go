package structutil

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shoraid/stx-go-utils/apperror"
	"github.com/stretchr/testify/assert"
)

func TestStructUtil_BindForm(t *testing.T) {
	type UserRequest struct {
		Name   string   `form:"name"`
		Age    int      `form:"age"`
		Active bool     `form:"active"`
		Tags   []string `form:"tags"`
	}

	tests := []struct {
		name     string
		formData url.Values
		expected UserRequest
		hasError bool
	}{
		{
			name: "Valid form with all fields",
			formData: url.Values{
				"name":   {"John Doe"},
				"age":    {"25"},
				"active": {"true"},
				"tags":   {"tag1", "tag2"},
			},
			expected: UserRequest{
				Name:   "John Doe",
				Age:    25,
				Active: true,
				Tags:   []string{"tag1", "tag2"},
			},
			hasError: false,
		},
		{
			name: "Partial form data",
			formData: url.Values{
				"name": {"Alice"},
			},
			expected: UserRequest{
				Name:   "Alice",
				Age:    0,
				Active: false,
				Tags:   nil,
			},
			hasError: false,
		},
		{
			name:     "Empty form data",
			formData: url.Values{},
			expected: UserRequest{
				Name:   "",
				Age:    0,
				Active: false,
				Tags:   nil,
			},
			hasError: false,
		},
		{
			name: "Form with false boolean",
			formData: url.Values{
				"name":   {"Bob"},
				"active": {"false"},
			},
			expected: UserRequest{
				Name:   "Bob",
				Active: false,
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			var result UserRequest
			err := BindForm(req, &result)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestStructUtil_BindForm_TypeConversion(t *testing.T) {
	type TypeRequest struct {
		IntField    int     `form:"int_field"`
		Int64Field  int64   `form:"int64_field"`
		Float64     float64 `form:"float64_field"`
		UintField   uint    `form:"uint_field"`
		StringField string  `form:"string_field"`
		BoolField   bool    `form:"bool_field"`
	}

	tests := []struct {
		name     string
		formData url.Values
		expected TypeRequest
		hasError bool
	}{
		{
			name: "Valid type conversions",
			formData: url.Values{
				"int_field":     {"-42"},
				"int64_field":   {"9223372036854775807"},
				"float64_field": {"3.14159"},
				"uint_field":    {"100"},
				"string_field":  {"hello world"},
				"bool_field":    {"1"},
			},
			expected: TypeRequest{
				IntField:    -42,
				Int64Field:  9223372036854775807,
				Float64:     3.14159,
				UintField:   100,
				StringField: "hello world",
				BoolField:   true,
			},
			hasError: false,
		},
		{
			name: "Boolean variations",
			formData: url.Values{
				"bool_field": {"true"},
			},
			expected: TypeRequest{
				BoolField: true,
			},
			hasError: false,
		},
		{
			name: "Invalid int value",
			formData: url.Values{
				"int_field": {"not-a-number"},
			},
			hasError: true,
		},
		{
			name: "Invalid float value",
			formData: url.Values{
				"float64_field": {"not-a-float"},
			},
			hasError: true,
		},
		{
			name: "Invalid bool value",
			formData: url.Values{
				"bool_field": {"not-a-bool"},
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			var result TypeRequest
			err := BindForm(req, &result)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestStructUtil_BindForm_PointerFields(t *testing.T) {
	type PointerRequest struct {
		Name   *string `form:"name"`
		Age    *int    `form:"age"`
		Active *bool   `form:"active"`
	}

	t.Run("Pointer fields with values", func(t *testing.T) {
		formData := url.Values{
			"name":   {"Alice"},
			"age":    {"30"},
			"active": {"true"},
		}

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var result PointerRequest
		err := BindForm(req, &result)

		assert.NoError(t, err)
		assert.NotNil(t, result.Name)
		assert.Equal(t, "Alice", *result.Name)
		assert.NotNil(t, result.Age)
		assert.Equal(t, 30, *result.Age)
		assert.NotNil(t, result.Active)
		assert.Equal(t, true, *result.Active)
	})

	t.Run("Pointer fields with empty values remain nil", func(t *testing.T) {
		formData := url.Values{
			"name": {""},
		}

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var result PointerRequest
		err := BindForm(req, &result)

		assert.NoError(t, err)
		assert.Nil(t, result.Name)
		assert.Nil(t, result.Age)
		assert.Nil(t, result.Active)
	})
}

func TestStructUtil_ValidateForm(t *testing.T) {
	type UserRequest struct {
		Name     string `form:"name" validate:"required,max=10"`
		Email    string `form:"email" validate:"required,email"`
		Age      int    `form:"age" validate:"min=18"`
		IsActive bool   `form:"is_active" validate:"boolean"`
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateForm(tt.request)

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

func TestStructUtil_ValidateForm_Nested(t *testing.T) {
	type Meta struct {
		Note string `form:"note" validate:"required"`
	}

	type Role struct {
		ID   string `form:"id" validate:"required,uuid"`
		Name string `form:"name" validate:"required"`
	}

	type UserRequest struct {
		Meta          Meta     `form:"meta" validate:"required"`
		Roles         []Role   `form:"roles" validate:"required,dive"`
		PermissionIDs []string `form:"permission_ids" validate:"required,dive,uuid"`
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
				Roles: []Role{
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
				Roles:         []Role{{ID: "d290f1ee-6c54-4b01-90e6-d701748f0851", Name: "Admin"}},
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
				Roles: []Role{
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
				Roles:         []Role{{ID: "d290f1ee-6c54-4b01-90e6-d701748f0851", Name: "Admin"}},
			},
			expected: map[string][]string{
				"permission_ids.0": {"field must be a valid UUID"},
			},
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateForm(tt.request)

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

func TestStructUtil_BindAndValidateForm(t *testing.T) {
	type LoginRequest struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required,min=6"`
	}

	tests := []struct {
		name           string
		formData       url.Values
		expectedError  error
		expectedFields map[string][]string
	}{
		{
			name: "Valid form and valid data",
			formData: url.Values{
				"email":    {"test@example.com"},
				"password": {"secret123"},
			},
		},
		{
			name: "Valid form but invalid data",
			formData: url.Values{
				"email":    {""},
				"password": {"123"},
			},
			expectedError: apperror.Err400InvalidData,
			expectedFields: map[string][]string{
				"email":    {"field is required"},
				"password": {"minimum value is 6"},
			},
		},
		{
			name:          "Empty form data",
			formData:      url.Values{},
			expectedError: apperror.Err400InvalidData,
			expectedFields: map[string][]string{
				"email":    {"field is required"},
				"password": {"field is required"},
			},
		},
		{
			name: "Invalid email format",
			formData: url.Values{
				"email":    {"not-an-email"},
				"password": {"secret123"},
			},
			expectedError: apperror.Err400InvalidData,
			expectedFields: map[string][]string{
				"email": {"field must be a valid email address"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			var input LoginRequest
			result, err := BindAndValidateForm(req, &input)

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

func TestStructUtil_getFormTagName(t *testing.T) {
	type TestStruct struct {
		WithTag      string `form:"with_tag"`
		WithTagOmit  string `form:"with_tag_omit,omitempty"`
		WithoutTag   string
		IgnoredField string `form:"-"`
		EmptyTag     string `form:""`
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

			result := getFormTagName(field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStructUtil_FormTypeError(t *testing.T) {
	err := &FormTypeError{
		Field:    "age",
		Expected: "int",
		Got:      "not-a-number",
	}

	assert.Equal(t, "field age: cannot convert 'not-a-number' to int", err.Error())
}

func BenchmarkStructutil_BindForm(b *testing.B) {
	type UserRequest struct {
		Name   string   `form:"name"`
		Age    int      `form:"age"`
		Active bool     `form:"active"`
		Tags   []string `form:"tags"`
	}

	formData := url.Values{
		"name":   {"John Doe"},
		"age":    {"25"},
		"active": {"true"},
		"tags":   {"tag1", "tag2", "tag3"},
	}

	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var result UserRequest
		BindForm(req, &result)
	}
}

func BenchmarkStructutil_ValidateForm(b *testing.B) {
	type UserRequest struct {
		Name     string `form:"name" validate:"required,max=10"`
		Email    string `form:"email" validate:"required,email"`
		Age      int    `form:"age" validate:"min=18"`
		IsActive bool   `form:"is_active" validate:"boolean"`
	}

	valid := UserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Age:      25,
		IsActive: true,
	}

	invalid := UserRequest{
		Name:     "",
		Email:    "invalidemail",
		Age:      10,
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
				ValidateForm(tt.payload)
			}
		})
	}
}

func BenchmarkStructutil_BindAndValidateForm(b *testing.B) {
	type LoginRequest struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required,min=6"`
	}

	tests := []struct {
		name     string
		formData url.Values
	}{
		{
			"Valid form and valid data",
			url.Values{"email": {"test@example.com"}, "password": {"secret123"}},
		},
		{
			"Valid form but invalid data",
			url.Values{"email": {""}, "password": {"123"}},
		},
		{
			"Empty form data",
			url.Values{},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.formData.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

				var input LoginRequest
				BindAndValidateForm(req, &input)
			}
		})
	}
}

func BenchmarkStructutil_getFormTagName(b *testing.B) {
	type BenchmarkStruct struct {
		WithTag      string `form:"with_tag"`
		WithTagOmit  string `form:"with_tag_omit,omitempty"`
		WithoutTag   string
		IgnoredField string `form:"-"`
		EmptyTag     string `form:""`
	}

	tType := reflect.TypeOf(BenchmarkStruct{})
	fields := []reflect.StructField{
		tType.Field(0),
		tType.Field(1),
		tType.Field(2),
		tType.Field(3),
		tType.Field(4),
	}

	for b.Loop() {
		for _, field := range fields {
			getFormTagName(field)
		}
	}
}

func TestStructUtil_BindForm_InvalidInput(t *testing.T) {
	t.Run("Nil pointer input", func(t *testing.T) {
		formData := url.Values{"name": {"test"}}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		var nilPtr *struct{ Name string }
		err := BindForm(req, nilPtr)
		assert.Error(t, err)
	})

	t.Run("Non-pointer input", func(t *testing.T) {
		formData := url.Values{"name": {"test"}}
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		input := struct{ Name string }{}
		err := BindForm(req, input)
		assert.Error(t, err)
	})
}

func TestStructUtil_BindForm_FieldWithNoFormTag(t *testing.T) {
	type Request struct {
		WithTag    string `form:"with_tag"`
		WithoutTag string
	}

	formData := url.Values{
		"with_tag":   {"value1"},
		"WithoutTag": {"value2"},
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var result Request
	err := BindForm(req, &result)

	assert.NoError(t, err)
	assert.Equal(t, "value1", result.WithTag)
	assert.Equal(t, "", result.WithoutTag) // Fields without form tag are not bound
}

func TestStructUtil_ValidateForm_FieldWithNoFormTag(t *testing.T) {
	type Request struct {
		NoTag string `validate:"required"`
	}

	result, err := ValidateForm(Request{NoTag: ""})

	assert.Equal(t, apperror.Err400InvalidData, err)
	assert.Equal(t, map[string][]string{
		"NoTag": {"field is required"},
	}, result)
}

func TestStructUtil_ValidateForm_FieldWithIgnoredFormTag(t *testing.T) {
	type Request struct {
		Ignored string `form:"-" validate:"required"`
	}

	result, err := ValidateForm(Request{Ignored: ""})

	assert.Equal(t, apperror.Err400InvalidData, err)
	assert.Equal(t, map[string][]string{
		"Ignored": {"field is required"},
	}, result)
}

func TestStructUtil_getFormErrorMessage(t *testing.T) {
	t.Run("FormTypeError", func(t *testing.T) {
		err := &FormTypeError{
			Field:    "age",
			Expected: "int",
			Got:      "abc",
		}

		fieldErrors, appErr := getFormErrorMessage(err)

		assert.Equal(t, apperror.Err400InvalidData, appErr)
		assert.Equal(t, map[string][]string{
			"age": {"invalid type, expected int"},
		}, fieldErrors)
	})

	t.Run("Unknown error returns nil", func(t *testing.T) {
		err := apperror.Err500InternalServer

		fieldErrors, appErr := getFormErrorMessage(err)

		assert.Nil(t, appErr)
		assert.Nil(t, fieldErrors)
	})
}

func TestStructUtil_BindAndValidateForm_TypeConversionError(t *testing.T) {
	type Request struct {
		Age int `form:"age" validate:"required"`
	}

	formData := url.Values{
		"age": {"not-a-number"},
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var input Request
	result, err := BindAndValidateForm(req, &input)

	assert.Equal(t, apperror.Err400InvalidData, err)
	assert.Equal(t, map[string][]string{
		"age": {"invalid type, expected int"},
	}, result)
}

func TestStructUtil_BindForm_SliceOfInts(t *testing.T) {
	type Request struct {
		IDs []int `form:"ids"`
	}

	formData := url.Values{
		"ids": {"1", "2", "3"},
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var result Request
	err := BindForm(req, &result)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result.IDs)
}

func TestStructUtil_ValidateForm_getErrorMessage(t *testing.T) {
	type Sample struct {
		Email    string `form:"email" validate:"email"`
		Name     string `form:"name" validate:"required"`
		Age      int    `form:"age" validate:"min=18"`
		IsActive bool   `form:"is_active" validate:"boolean"`
		Role     string `form:"role" validate:"oneof=admin user"`
		ID       string `form:"id" validate:"uuid"`
		MaxTest  string `form:"max_test" validate:"max=5"`
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
		{"boolean", "field must be a boolean"},
		{"oneof", "field must be one of: admin, user"},
		{"uuid", "field must be a valid UUID"},
		{"max", "maximum length is 5"},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			var matched validator.FieldError
			for _, fe := range validationErrors {
				if fe.Tag() == tt.tag {
					matched = fe
					break
				}
			}
			if matched != nil {
				actual := getErrorMessage(matched)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

// createMultipartRequest creates a multipart form request with files and fields
func createMultipartRequest(fields map[string]string, files map[string][]struct {
	filename string
	content  []byte
}) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add regular form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	// Add file fields
	for fieldName, fileList := range files {
		for _, file := range fileList {
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+file.filename+`"`)
			h.Set("Content-Type", "application/octet-stream")

			part, err := writer.CreatePart(h)
			if err != nil {
				return nil, err
			}
			if _, err := part.Write(file.content); err != nil {
				return nil, err
			}
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func TestStructUtil_BindForm_SingleFile(t *testing.T) {
	type UploadRequest struct {
		Name   string                `form:"name"`
		Avatar *multipart.FileHeader `form:"avatar"`
	}

	t.Run("Single file upload", func(t *testing.T) {
		fields := map[string]string{"name": "John"}
		files := map[string][]struct {
			filename string
			content  []byte
		}{
			"avatar": {{filename: "avatar.png", content: []byte("fake image content")}},
		}

		req, err := createMultipartRequest(fields, files)
		assert.NoError(t, err)

		var result UploadRequest
		err = BindForm(req, &result)

		assert.NoError(t, err)
		assert.Equal(t, "John", result.Name)
		assert.NotNil(t, result.Avatar)
		assert.Equal(t, "avatar.png", result.Avatar.Filename)
		assert.Equal(t, int64(len("fake image content")), result.Avatar.Size)
	})

	t.Run("No file uploaded", func(t *testing.T) {
		fields := map[string]string{"name": "Jane"}
		files := map[string][]struct {
			filename string
			content  []byte
		}{}

		req, err := createMultipartRequest(fields, files)
		assert.NoError(t, err)

		var result UploadRequest
		err = BindForm(req, &result)

		assert.NoError(t, err)
		assert.Equal(t, "Jane", result.Name)
		assert.Nil(t, result.Avatar)
	})
}

func TestStructUtil_BindForm_MultipleFiles(t *testing.T) {
	type GalleryRequest struct {
		Title  string                  `form:"title"`
		Photos []*multipart.FileHeader `form:"photos"`
	}

	t.Run("Multiple files upload", func(t *testing.T) {
		fields := map[string]string{"title": "My Gallery"}
		files := map[string][]struct {
			filename string
			content  []byte
		}{
			"photos": {
				{filename: "photo1.jpg", content: []byte("photo 1 content")},
				{filename: "photo2.jpg", content: []byte("photo 2 content")},
				{filename: "photo3.jpg", content: []byte("photo 3 content")},
			},
		}

		req, err := createMultipartRequest(fields, files)
		assert.NoError(t, err)

		var result GalleryRequest
		err = BindForm(req, &result)

		assert.NoError(t, err)
		assert.Equal(t, "My Gallery", result.Title)
		assert.Len(t, result.Photos, 3)
		assert.Equal(t, "photo1.jpg", result.Photos[0].Filename)
		assert.Equal(t, "photo2.jpg", result.Photos[1].Filename)
		assert.Equal(t, "photo3.jpg", result.Photos[2].Filename)
	})

	t.Run("Empty photos array", func(t *testing.T) {
		fields := map[string]string{"title": "Empty Gallery"}
		files := map[string][]struct {
			filename string
			content  []byte
		}{}

		req, err := createMultipartRequest(fields, files)
		assert.NoError(t, err)

		var result GalleryRequest
		err = BindForm(req, &result)

		assert.NoError(t, err)
		assert.Equal(t, "Empty Gallery", result.Title)
		assert.Nil(t, result.Photos)
	})
}

func TestStructUtil_BindForm_MixedFieldsAndFiles(t *testing.T) {
	type ProfileRequest struct {
		Name     string                  `form:"name"`
		Email    string                  `form:"email"`
		Age      int                     `form:"age"`
		Avatar   *multipart.FileHeader   `form:"avatar"`
		Photos   []*multipart.FileHeader `form:"photos"`
		IsActive bool                    `form:"is_active"`
	}

	fields := map[string]string{
		"name":      "Alice",
		"email":     "alice@example.com",
		"age":       "25",
		"is_active": "true",
	}
	files := map[string][]struct {
		filename string
		content  []byte
	}{
		"avatar": {{filename: "profile.png", content: []byte("profile pic")}},
		"photos": {
			{filename: "vacation1.jpg", content: []byte("vacation 1")},
			{filename: "vacation2.jpg", content: []byte("vacation 2")},
		},
	}

	req, err := createMultipartRequest(fields, files)
	assert.NoError(t, err)

	var result ProfileRequest
	err = BindForm(req, &result)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", result.Name)
	assert.Equal(t, "alice@example.com", result.Email)
	assert.Equal(t, 25, result.Age)
	assert.True(t, result.IsActive)
	assert.NotNil(t, result.Avatar)
	assert.Equal(t, "profile.png", result.Avatar.Filename)
	assert.Len(t, result.Photos, 2)
}

func TestStructUtil_isFileField(t *testing.T) {
	fileHeaderType := reflect.TypeOf((*multipart.FileHeader)(nil))
	fileHeaderSliceType := reflect.TypeOf([]*multipart.FileHeader{})
	stringType := reflect.TypeOf("")
	intType := reflect.TypeOf(0)

	assert.True(t, isFileField(fileHeaderType))
	assert.True(t, isFileField(fileHeaderSliceType))
	assert.False(t, isFileField(stringType))
	assert.False(t, isFileField(intType))
}

func TestStructUtil_BindForm_FileWithValidation(t *testing.T) {
	type UploadRequest struct {
		Name   string                `form:"name" validate:"required"`
		Avatar *multipart.FileHeader `form:"avatar" validate:"required"`
	}

	t.Run("Missing required file", func(t *testing.T) {
		fields := map[string]string{"name": "Bob"}
		files := map[string][]struct {
			filename string
			content  []byte
		}{}

		req, err := createMultipartRequest(fields, files)
		assert.NoError(t, err)

		var result UploadRequest
		fieldErrors, err := BindAndValidateForm(req, &result)

		assert.Equal(t, apperror.Err400InvalidData, err)
		assert.Equal(t, map[string][]string{
			"avatar": {"field is required"},
		}, fieldErrors)
	})

	t.Run("Valid file and field", func(t *testing.T) {
		fields := map[string]string{"name": "Bob"}
		files := map[string][]struct {
			filename string
			content  []byte
		}{
			"avatar": {{filename: "bob.png", content: []byte("bob's picture")}},
		}

		req, err := createMultipartRequest(fields, files)
		assert.NoError(t, err)

		var result UploadRequest
		fieldErrors, err := BindAndValidateForm(req, &result)

		assert.NoError(t, err)
		assert.Nil(t, fieldErrors)
		assert.Equal(t, "Bob", result.Name)
		assert.NotNil(t, result.Avatar)
	})
}

func BenchmarkStructutil_BindForm_WithFiles(b *testing.B) {
	type UploadRequest struct {
		Name   string                  `form:"name"`
		Avatar *multipart.FileHeader   `form:"avatar"`
		Photos []*multipart.FileHeader `form:"photos"`
	}

	fields := map[string]string{"name": "Benchmark User"}
	files := map[string][]struct {
		filename string
		content  []byte
	}{
		"avatar": {{filename: "avatar.png", content: []byte("avatar content")}},
		"photos": {
			{filename: "photo1.jpg", content: []byte("photo 1")},
			{filename: "photo2.jpg", content: []byte("photo 2")},
		},
	}

	for b.Loop() {
		req, _ := createMultipartRequest(fields, files)
		var result UploadRequest
		BindForm(req, &result)
	}
}
