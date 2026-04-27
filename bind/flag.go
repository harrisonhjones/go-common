package bind

import (
	"context"
	"flag"
	"fmt"
	"reflect"
)

// FromFlags populates the struct pointed to by v with values parsed from
// command-line flags. args should typically be os.Args[1:]. v must be a
// non-nil pointer to a struct. See the package documentation for supported
// struct tags and field types.
func FromFlags(ctx context.Context, args []string, v any) error {
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

	fs := flag.NewFlagSet("bind", flag.ContinueOnError)

	// Register a string flag for each tagged field. We'll use set() for type conversion.
	flagPtrs := make(map[int]*string)
	for i := range vTypeElem.NumField() {
		typeField := vTypeElem.Field(i)

		name := typeField.Tag.Get("flag")
		if name == "" {
			continue
		}

		defaultVal := typeField.Tag.Get("flagDefault")
		usage := typeField.Tag.Get("flagUsage")

		ptr := fs.String(name, defaultVal, usage)
		flagPtrs[i] = ptr
	}

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// Track which flags were explicitly set on the command line.
	setFlags := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	vValElem := reflect.ValueOf(v).Elem()

	for i, ptr := range flagPtrs {
		typeField := vTypeElem.Field(i)
		valField := vValElem.Field(i)
		flagName := typeField.Tag.Get("flag")

		// If the flag wasn't explicitly provided and there's no default, skip.
		if !setFlags[flagName] && *ptr == "" {
			continue
		}

		if typeField.Tag.Get("flagFromNoOverwrite") == "true" && !valField.IsZero() {
			continue
		}

		if err := set(typeField.Type, valField, *ptr); err != nil {
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
