package structutil

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shoraid/stx-go-utils/apperror"
)

var validate = validator.New()

// validates any struct using struct tags
func Validate(input any) (map[string][]string, error) {
	err := validate.Struct(input)
	if err == nil {
		return nil, nil
	}

	validationErrors := make(map[string][]string)
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	if val.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for _, fe := range err.(validator.ValidationErrors) {
		fieldName := getJSONFieldName(fe.StructField(), typ)
		message := getErrorMessage(fe)
		validationErrors[fieldName] = append(validationErrors[fieldName], message)
	}

	return validationErrors, apperror.Err400InvalidData
}

func getJSONFieldName(structField string, typ reflect.Type) string {
	if field, ok := typ.FieldByName(structField); ok {
		tag := field.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name != "" && name != "-" {
			return name
		}
	}
	return structField
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "max":
		return "maximum length is " + fe.Param()
	case "min":
		return "minimum value is " + fe.Param()
	case "boolean":
		return "field must be a boolean"
	case "oneof":
		return "field must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	default:
		return "field is invalid"
	}
}
