package structutil

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shoraid/stx-go-utils/apperror"
)

var validate = validator.New()

// Validate validates a struct using `validate` tags and returns a map of field errors
// using JSON tag names. Supports nested structs and slices.
//
// Parameters:
// - input: struct or pointer to struct with `validate` tags.
//
// Returns:
// - map[string][]string: validation errors using JSON field paths as keys.
// - error: apperror.Err400InvalidData if validation fails, nil if valid.
//
// Features:
// - Uses reflection to get JSON tag names for error keys.
// - Supports flat fields, nested fields, and slice elements (with index).
//
// Example:
//
//	type Role struct {
//	    ID   string `json:"id" validate:"required,uuid"`
//	    Name string `json:"name" validate:"required"`
//	}
//
//	type Meta struct {
//	    Note string `json:"note" validate:"required"`
//	}
//
//	type UserRequest struct {
//	    Name          string   `json:"name" validate:"required"`
//	    Meta          Meta     `json:"meta"`
//	    Roles         []Role   `json:"roles" validate:"dive"`
//	    PermissionIDs []string `json:"permissionIds" validate:"dive,uuid"`
//	}
//
//	input := UserRequest{
//	    Name: "", // invalid
//	    Meta: Meta{Note: ""}, // invalid
//	    Roles: []Role{
//	        {ID: "invalid-uuid", Name: ""}, // both invalid
//	    },
//	    PermissionIDs: []string{"invalid-uuid"}, // invalid
//	}
//
//	Validate(input)
//	// Output:
//	map[string][]string{
//	    "name":              {"field is required"},
//	    "meta.note":         {"field is required"},
//	    "roles.0.id":        {"field must be a valid UUID"},
//	    "roles.0.name":      {"field is required"},
//	    "permissionIds.0":   {"field must be a valid UUID"},
//	}, apperror.Err400InvalidData
func Validate(input any) (map[string][]string, error) {
	err := validate.Struct(input)
	if err == nil {
		return nil, nil
	}

	validationErrors := make(map[string][]string)

	root := reflect.TypeOf(input)
	if root.Kind() == reflect.Ptr {
		root = root.Elem()
	}

	for _, fe := range err.(validator.ValidationErrors) {
		fieldPath := buildJSONPath(root, fe)
		message := getErrorMessage(fe)
		validationErrors[fieldPath] = append(validationErrors[fieldPath], message)
	}

	return validationErrors, apperror.Err400InvalidData
}

func BindAndValidateJSON(r *http.Request, input any) (map[string][]string, error) {
	err := BindJSON(r, input)
	if err != nil {

		fieldErrors, jsonErr := getJsonErrorMessage(err)
		if jsonErr != nil {
			return fieldErrors, jsonErr
		}
	}

	return Validate(input)
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "field must be a valid email address"
	case "max":
		return "maximum length is " + fe.Param()
	case "min":
		return "minimum value is " + fe.Param()
	case "boolean":
		return "field must be a boolean"
	case "oneof":
		return "field must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	case "uuid":
		return "field must be a valid UUID"
	default:
		return "field is invalid"
	}
}

func getJSONTagName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	name := strings.Split(tag, ",")[0]
	if name == "" || name == "-" {
		return field.Name
	}
	return name
}

func buildJSONPath(root reflect.Type, fe validator.FieldError) string {
	ns := fe.StructNamespace() // e.g. "Meta.Note" or "Items[0].Name"
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
				jsonKey := getJSONTagName(field)
				path = append(path, jsonKey+"."+index)

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
			jsonKey := getJSONTagName(field)
			path = append(path, jsonKey)

			current = field.Type
			if current.Kind() == reflect.Ptr {
				current = current.Elem()
			}
		}
	}

	return strings.Join(path, ".")
}

func getJsonErrorMessage(err error) (map[string][]string, error) {
	switch e := err.(type) {
	case *json.SyntaxError:
		return map[string][]string{
			"json": {"invalid JSON format: please check for missing commas, braces, or quotes"},
		}, apperror.Err400InvalidBody
	case *json.UnmarshalTypeError:
		fieldName := e.Field
		if fieldName != "" {
			return map[string][]string{
				fieldName: {"invalid type, expected " + e.Type.String()},
			}, apperror.Err400InvalidData
		}
	}

	return nil, nil
}
