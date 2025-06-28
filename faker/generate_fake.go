package faker

import (
	"reflect"
	"time"
)

func GenerateFake[T any]() *T {
	t := new(T)
	v := reflect.ValueOf(t).Elem()
	tType := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := tType.Field(i)
		tag := fieldType.Tag.Get("faker")

		if !field.CanSet() || tag == "" {
			continue
		}

		switch tag {
		case "string":
			val := RandString(20)
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&val))
			} else {
				field.SetString(val)
			}
		case "sentence":
			val := RandSentence(2)
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&val))
			} else {
				field.SetString(val)
			}
		case "uuid_str":
			val := UUID()
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&val))
			} else {
				field.Set(reflect.ValueOf(val))
			}
		case "bool":
			val := RandBool()
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&val))
			} else {
				field.SetBool(val)
			}
		case "int":
			val := RandInt(99, 9999)
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&val))
			} else {
				field.SetInt(int64(val))
			}
		case "time":
			val := RandTime(time.Now(), time.Now().Add(30*24*time.Hour))
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.ValueOf(&val))
			} else {
				field.Set(reflect.ValueOf(val))
			}
		}
	}

	return t
}
