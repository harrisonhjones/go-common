package bind

import (
	"context"
	"fmt"
	"os"
	"reflect"
)

// FromEnv loads data into v with values from the local environment.
func FromEnv(ctx context.Context, v any) error {
	if v == nil {
		return fmt.Errorf("v must not be nil")
	}

	vType := reflect.TypeOf(v)
	if vType.Kind() != reflect.Pointer {
		return fmt.Errorf("v must be a pointer")
	}

	vTypeElem := vType.Elem()
	if vTypeElem.Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to a struct")
	}

	vValElem := reflect.ValueOf(v).Elem()

	for i := range vTypeElem.NumField() {
		typeField := vTypeElem.Field(i)

		var val string
		var key string
		if key = typeField.Tag.Get("env"); key != "" {
			val = os.Getenv(key)
		}
		if val == "" {
			defaultVal := typeField.Tag.Get("envDefault")
			if defaultVal == "" {
				continue
			}
			val = defaultVal
		}

		valField := vValElem.Field(i)

		// If envFromNoOverwrite is set, skip fields that already have a non-zero value.
		if typeField.Tag.Get("envFromNoOverwrite") == "true" && !valField.IsZero() {
			continue
		}

		if err := set(typeField.Type, valField, val); err != nil {
			return err
		}
	}

	if normalizingVal, ok := v.(Normalizer); ok {
		normalizingVal.Normalize(ctx)
	}

	if validatingVal, ok := v.(Validator); ok {
		if err := validatingVal.Validate(ctx); err != nil {
			return fmt.Errorf("validation failed: %v", err)
		}
	}

	return nil
}
