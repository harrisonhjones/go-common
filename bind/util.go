package bind

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Normalizer is implemented by structs that need to adjust field values after
// binding (e.g. lowercasing hostnames). Normalize is called after all fields
// are set but before validation.
type Normalizer interface {
	Normalize(context.Context)
}

// Validator is implemented by structs that need to check invariants after
// binding. Validate is called after normalization; a non-nil return value
// causes the bind function to return an error.
type Validator interface {
	Validate(context.Context) error
}

func set(fieldType reflect.Type, fieldValue reflect.Value, value string) error {
	switch fieldKind := fieldType.Kind(); fieldKind {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Bool:
		fieldValue.SetBool(strToBool(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return fmt.Errorf("failed to parse field as integer: %#v", err)
		}
		fieldValue.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return fmt.Errorf("failed to parse field as uint: %#v", err)
		}
		fieldValue.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse field as float: %#v", err)
		}
		fieldValue.SetFloat(i)
	case reflect.Slice:
		switch elemKind := fieldType.Elem().Kind(); elemKind {
		case reflect.String:
			fieldValue.Set(reflect.ValueOf(strings.Split(value, ",")))
		default:
			return fmt.Errorf("slice field type %q not supported; unable to bind value", elemKind.String())
		}
	default:
		return fmt.Errorf("field type %q not supported; unable to bind value", fieldKind.String())
	}
	return nil
}

func strToBool(val string) bool {
	if v := strings.ToLower(val); v == "y" || v == "yes" || v == "true" {
		return true
	}
	return false
}
