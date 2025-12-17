package structutil

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shoraid/stx-go-utils/apperror"
)

// BindForm binds form data from an HTTP request to a struct using the `form` tag.
// Supports both application/x-www-form-urlencoded and multipart/form-data.
//
// Parameters:
// - r: HTTP request with form data.
// - input: pointer to struct with `form` tags.
//
// Returns:
// - error: binding error if form parsing fails or type conversion fails.
//
// Supported field types:
// - Scalar: string, int, int64, float64, bool, uint and their pointer variants
// - Slices: []string, []int, etc.
// - Files: *multipart.FileHeader (single file), []*multipart.FileHeader (multiple files)
//
// Example:
//
//	type CreateUserRequest struct {
//	    Name    string                  `form:"name"`
//	    Age     int                     `form:"age"`
//	    Active  bool                    `form:"active"`
//	    Tags    []string                `form:"tags"`
//	    Avatar  *multipart.FileHeader   `form:"avatar"`   // Single file
//	    Photos  []*multipart.FileHeader `form:"photos"`   // Multiple files
//	}
//
//	var input CreateUserRequest
//	err := BindForm(r, &input)
func BindForm(r *http.Request, input any) error {
	contentType := r.Header.Get("Content-Type")

	var multipartForm *multipart.Form

	// Parse the form based on content type
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MB max memory
			return err
		}
		multipartForm = r.MultipartForm
	} else {
		if err := r.ParseForm(); err != nil {
			return err
		}
	}

	return bindFormValues(r.Form, multipartForm, input)
}

// bindFormValues binds url.Values and file uploads to a struct using reflection
func bindFormValues(values map[string][]string, multipartForm *multipart.Form, input any) error {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return apperror.Err400InvalidBody
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return apperror.Err400InvalidBody
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		formTag := field.Tag.Get("form")
		if formTag == "" || formTag == "-" {
			continue
		}

		formKey := strings.Split(formTag, ",")[0]

		// Check if this is a file field
		if isFileField(field.Type) {
			if multipartForm != nil {
				if err := setFileFieldValue(fieldValue, field.Type, multipartForm.File[formKey]); err != nil {
					return err
				}
			}
			continue
		}

		// Handle regular form values
		formValues, exists := values[formKey]
		if !exists || len(formValues) == 0 {
			continue
		}

		if err := setFieldValue(fieldValue, formValues); err != nil {
			return &FormTypeError{
				Field:    formKey,
				Expected: field.Type.String(),
				Got:      formValues[0],
			}
		}
	}

	return nil
}

// isFileField checks if the field type is a file-related type
func isFileField(t reflect.Type) bool {
	fileHeaderType := reflect.TypeOf((*multipart.FileHeader)(nil))
	fileHeaderSliceType := reflect.TypeOf([]*multipart.FileHeader{})

	return t == fileHeaderType || t == fileHeaderSliceType
}

// setFileFieldValue sets file field values from multipart form
func setFileFieldValue(fieldValue reflect.Value, fieldType reflect.Type, files []*multipart.FileHeader) error {
	if len(files) == 0 {
		return nil
	}

	fileHeaderType := reflect.TypeOf((*multipart.FileHeader)(nil))
	fileHeaderSliceType := reflect.TypeOf([]*multipart.FileHeader{})

	switch fieldType {
	case fileHeaderType:
		// Single file: *multipart.FileHeader
		fieldValue.Set(reflect.ValueOf(files[0]))
	case fileHeaderSliceType:
		// Multiple files: []*multipart.FileHeader
		fieldValue.Set(reflect.ValueOf(files))
	}

	return nil
}

// FormTypeError represents a type conversion error during form binding
type FormTypeError struct {
	Field    string
	Expected string
	Got      string
}

func (e *FormTypeError) Error() string {
	return "field " + e.Field + ": cannot convert '" + e.Got + "' to " + e.Expected
}

// setFieldValue sets the value of a struct field based on form values
func setFieldValue(fieldValue reflect.Value, values []string) error {
	fieldType := fieldValue.Type()

	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		if values[0] == "" {
			return nil // Leave nil for empty values
		}
		// Create a new value and set it
		newValue := reflect.New(fieldType.Elem())
		if err := setFieldValue(newValue.Elem(), values); err != nil {
			return err
		}
		fieldValue.Set(newValue)
		return nil
	}

	// Handle slice types (except string slices handled specially)
	if fieldType.Kind() == reflect.Slice {
		if fieldType.Elem().Kind() == reflect.String {
			fieldValue.Set(reflect.ValueOf(values))
			return nil
		}
		// For other slice types, try to convert each value
		slice := reflect.MakeSlice(fieldType, len(values), len(values))
		for i, v := range values {
			if err := setScalarValue(slice.Index(i), v); err != nil {
				return err
			}
		}
		fieldValue.Set(slice)
		return nil
	}

	return setScalarValue(fieldValue, values[0])
}

// setScalarValue sets a scalar value from a string
func setScalarValue(fieldValue reflect.Value, value string) error {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			return nil
		}
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			return nil
		}
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		if value == "" {
			return nil
		}
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(floatVal)

	case reflect.Bool:
		if value == "" {
			return nil
		}
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		fieldValue.SetBool(boolVal)

	default:
		return &FormTypeError{
			Expected: fieldValue.Type().String(),
			Got:      value,
		}
	}

	return nil
}

// ValidateForm validates a struct using `validate` tags and returns a map of field errors
// using form tag names. Supports nested structs and slices.
//
// Parameters:
// - input: struct or pointer to struct with `validate` tags.
//
// Returns:
// - map[string][]string: validation errors using form field paths as keys.
// - error: apperror.Err400InvalidData if validation fails, nil if valid.
//
// Example:
//
//	type CreateUserRequest struct {
//	    Name  string `form:"name" validate:"required,max=100"`
//	    Email string `form:"email" validate:"required,email"`
//	    Age   int    `form:"age" validate:"min=18"`
//	}
//
//	ValidateForm(input)
//	// Output:
//	map[string][]string{
//	    "name":  {"field is required"},
//	    "email": {"field must be a valid email address"},
//	    "age":   {"minimum value is 18"},
//	}, apperror.Err400InvalidData
func ValidateForm(input any) (map[string][]string, error) {
	err := Validator.Struct(input)
	if err == nil {
		return nil, nil
	}

	validationErrors := make(map[string][]string)

	root := reflect.TypeOf(input)
	if root.Kind() == reflect.Pointer {
		root = root.Elem()
	}

	for _, fe := range err.(validator.ValidationErrors) {
		fieldPath := buildFormPath(root, fe)
		message := getErrorMessage(fe)
		validationErrors[fieldPath] = append(validationErrors[fieldPath], message)
	}

	return validationErrors, apperror.Err400InvalidData
}

// BindAndValidateForm binds form data to a struct and validates it.
//
// Parameters:
// - r: HTTP request with form data.
// - input: pointer to struct with `form` and `validate` tags.
//
// Returns:
// - map[string][]string: validation errors using form field names as keys.
// - error: apperror.Err400InvalidBody if binding fails, apperror.Err400InvalidData if validation fails.
//
// Example:
//
//	type CreateUserRequest struct {
//	    Name  string `form:"name" validate:"required"`
//	    Email string `form:"email" validate:"required,email"`
//	}
//
//	var input CreateUserRequest
//	fieldErrors, err := BindAndValidateForm(r, &input)
func BindAndValidateForm(r *http.Request, input any) (map[string][]string, error) {
	err := BindForm(r, input)
	if err != nil {
		fieldErrors, formErr := getFormErrorMessage(err)
		if formErr != nil {
			return fieldErrors, formErr
		}
	}

	return ValidateForm(input)
}

// getFormTagName returns the form tag name or falls back to the field name
func getFormTagName(field reflect.StructField) string {
	tag := field.Tag.Get("form")
	name := strings.Split(tag, ",")[0]
	if name == "" || name == "-" {
		return field.Name
	}
	return name
}

// buildFormPath builds the form field path from validation error
func buildFormPath(root reflect.Type, fe validator.FieldError) string {
	ns := fe.StructNamespace()
	parts := strings.Split(ns, ".")

	var path []string
	current := root

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Handle index (slice), e.g. Items[0]
		if strings.Contains(part, "[") {
			name := part[:strings.Index(part, "[")]
			index := part[strings.Index(part, "[")+1 : strings.Index(part, "]")]

			if field, ok := current.FieldByName(name); ok {
				formKey := getFormTagName(field)
				path = append(path, formKey+"."+index)

				current = field.Type
				if current.Kind() == reflect.Slice {
					current = current.Elem()
				}
				if current.Kind() == reflect.Ptr {
					current = current.Elem()
				}
			}
			continue
		}

		if field, ok := current.FieldByName(part); ok {
			formKey := getFormTagName(field)
			path = append(path, formKey)

			current = field.Type
			if current.Kind() == reflect.Ptr {
				current = current.Elem()
			}
		}
	}

	return strings.Join(path, ".")
}

// getFormErrorMessage converts binding errors to field error maps
func getFormErrorMessage(err error) (map[string][]string, error) {
	switch e := err.(type) {
	case *FormTypeError:
		return map[string][]string{
			e.Field: {"invalid type, expected " + e.Expected},
		}, apperror.Err400InvalidData
	case *json.SyntaxError:
		return map[string][]string{
			"form": {"invalid form data format"},
		}, apperror.Err400InvalidBody
	}

	return nil, nil
}
