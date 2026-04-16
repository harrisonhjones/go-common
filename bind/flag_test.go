package bind

import (
	"context"
	"fmt"
	"testing"
)

// --- test structs ---

type flagBasicStruct struct {
	Host       string   `flag:"host"`
	Port       int      `flag:"port"`
	Verbose    bool     `flag:"verbose"`
	Rate       float64  `flag:"rate"`
	Count      uint     `flag:"count"`
	Tags       []string `flag:"tags"`
	Int8Val    int8     `flag:"int8val"`
	Int16Val   int16    `flag:"int16val"`
	Int32Val   int32    `flag:"int32val"`
	Int64Val   int64    `flag:"int64val"`
	Uint8Val   uint8    `flag:"uint8val"`
	Uint16Val  uint16   `flag:"uint16val"`
	Uint32Val  uint32   `flag:"uint32val"`
	Uint64Val  uint64   `flag:"uint64val"`
	Float32Val float32  `flag:"float32val"`
}

type flagDefaultStruct struct {
	Host string  `flag:"host" flagDefault:"localhost"`
	Port int     `flag:"port" flagDefault:"8080"`
	Rate float64 `flag:"rate" flagDefault:"1.5"`
}

type flagUsageStruct struct {
	Host string `flag:"host" flagDefault:"localhost" flagUsage:"server hostname"`
}

type flagNoTagStruct struct {
	Name string
	Age  int
}

type flagMixedStruct struct {
	Tagged   string `flag:"tagged"`
	Untagged string
}

type flagNormalizerStruct struct {
	Name       string `flag:"name"`
	normalized bool
}

func (n *flagNormalizerStruct) Normalize(_ context.Context) {
	n.normalized = true
	n.Name = "normalized:" + n.Name
}

type flagValidatorStruct struct {
	Name string `flag:"name"`
}

func (v *flagValidatorStruct) Validate(_ context.Context) error {
	if v.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type flagNoOverwriteStruct struct {
	Host string `flag:"host" flagFromNoOverwrite:"true"`
	Port int    `flag:"port" flagFromNoOverwrite:"true"`
}

// --- tests ---

func TestFromFlags_NonPointer(t *testing.T) {
	var s flagBasicStruct
	err := FromFlags(context.Background(), nil, s)
	if err == nil {
		t.Fatal("expected error for non-pointer, got nil")
	}
}

func TestFromFlags_PointerToNonStruct(t *testing.T) {
	s := "hello"
	err := FromFlags(context.Background(), nil, &s)
	if err == nil {
		t.Fatal("expected error for pointer to non-struct, got nil")
	}
}

func TestFromFlags_Nil(t *testing.T) {
	err := FromFlags(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error for nil, got nil")
	}
}

func TestFromFlags_BasicTypes(t *testing.T) {
	args := []string{
		"-host", "example.com",
		"-port", "-42",
		"-verbose", "true",
		"-rate", "3.14",
		"-count", "10",
		"-tags", "a,b,c",
		"-int8val", "127",
		"-int16val", "-300",
		"-int32val", "100000",
		"-int64val", "-9999999999",
		"-uint8val", "255",
		"-uint16val", "65535",
		"-uint32val", "4294967295",
		"-uint64val", "18446744073709551615",
		"-float32val", "2.5",
	}

	var s flagBasicStruct
	if err := FromFlags(context.Background(), args, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "Host", s.Host, "example.com")
	assertEqual(t, "Port", s.Port, -42)
	assertEqual(t, "Verbose", s.Verbose, true)
	assertFloat(t, "Rate", s.Rate, 3.14, 0.001)
	assertEqual(t, "Count", s.Count, uint(10))
	assertSliceEqual(t, "Tags", s.Tags, []string{"a", "b", "c"})
	assertEqual(t, "Int8Val", s.Int8Val, int8(127))
	assertEqual(t, "Int16Val", s.Int16Val, int16(-300))
	assertEqual(t, "Int32Val", s.Int32Val, int32(100000))
	assertEqual(t, "Int64Val", s.Int64Val, int64(-9999999999))
	assertEqual(t, "Uint8Val", s.Uint8Val, uint8(255))
	assertEqual(t, "Uint16Val", s.Uint16Val, uint16(65535))
	assertEqual(t, "Uint32Val", s.Uint32Val, uint32(4294967295))
	assertEqual(t, "Uint64Val", s.Uint64Val, uint64(18446744073709551615))
	assertFloat(t, "Float32Val", float64(s.Float32Val), 2.5, 0.01)
}

func TestFromFlags_Defaults(t *testing.T) {
	// No args provided — defaults should be used
	var s flagDefaultStruct
	if err := FromFlags(context.Background(), nil, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "Host", s.Host, "localhost")
	assertEqual(t, "Port", s.Port, 8080)
	assertFloat(t, "Rate", s.Rate, 1.5, 0.01)
}

func TestFromFlags_ArgsOverrideDefaults(t *testing.T) {
	args := []string{"-host", "override.com", "-port", "9090"}

	var s flagDefaultStruct
	if err := FromFlags(context.Background(), args, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "Host", s.Host, "override.com")
	assertEqual(t, "Port", s.Port, 9090)
	// Rate not provided — should use default
	assertFloat(t, "Rate", s.Rate, 1.5, 0.01)
}

func TestFromFlags_NoTags(t *testing.T) {
	var s flagNoTagStruct
	if err := FromFlags(context.Background(), nil, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Name", s.Name, "")
	assertEqual(t, "Age", s.Age, 0)
}

func TestFromFlags_MixedTags(t *testing.T) {
	args := []string{"-tagged", "value"}

	var s flagMixedStruct
	if err := FromFlags(context.Background(), args, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Tagged", s.Tagged, "value")
	assertEqual(t, "Untagged", s.Untagged, "")
}

func TestFromFlags_EmptyStruct(t *testing.T) {
	type empty struct{}
	var s empty
	if err := FromFlags(context.Background(), nil, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromFlags_EmptyArgs(t *testing.T) {
	var s flagBasicStruct
	if err := FromFlags(context.Background(), []string{}, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// All fields should be zero-valued
	assertEqual(t, "Host", s.Host, "")
	assertEqual(t, "Port", s.Port, 0)
	assertEqual(t, "Verbose", s.Verbose, false)
}

func TestFromFlags_InvalidFlag(t *testing.T) {
	args := []string{"-unknown", "value"}

	var s flagBasicStruct
	err := FromFlags(context.Background(), args, &s)
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestFromFlags_InvalidIntValue(t *testing.T) {
	args := []string{"-port", "not_a_number"}

	var s flagBasicStruct
	err := FromFlags(context.Background(), args, &s)
	if err == nil {
		t.Fatal("expected error for invalid int, got nil")
	}
}

func TestFromFlags_Normalizer(t *testing.T) {
	args := []string{"-name", "alice"}

	var s flagNormalizerStruct
	if err := FromFlags(context.Background(), args, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.normalized {
		t.Fatal("expected Normalize to be called")
	}
	assertEqual(t, "Name", s.Name, "normalized:alice")
}

func TestFromFlags_ValidatorPass(t *testing.T) {
	args := []string{"-name", "bob"}

	var s flagValidatorStruct
	if err := FromFlags(context.Background(), args, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromFlags_ValidatorFail(t *testing.T) {
	var s flagValidatorStruct
	err := FromFlags(context.Background(), nil, &s)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestFromFlags_NoOverwrite_Preserved(t *testing.T) {
	args := []string{"-host", "from_flag", "-port", "9090"}

	s := flagNoOverwriteStruct{Host: "already_set", Port: 42}
	if err := FromFlags(context.Background(), args, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Host", s.Host, "already_set")
	assertEqual(t, "Port", s.Port, 42)
}

func TestFromFlags_NoOverwrite_ZeroLoadsFlag(t *testing.T) {
	args := []string{"-host", "from_flag", "-port", "9090"}

	var s flagNoOverwriteStruct
	if err := FromFlags(context.Background(), args, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Host", s.Host, "from_flag")
	assertEqual(t, "Port", s.Port, 9090)
}

func TestFromFlags_BoolVariants(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"true", []string{"-verbose", "true"}, true},
		{"True", []string{"-verbose", "True"}, true},
		{"y", []string{"-verbose", "y"}, true},
		{"yes", []string{"-verbose", "yes"}, true},
		{"YES", []string{"-verbose", "YES"}, true},
		{"false", []string{"-verbose", "false"}, false},
		{"no", []string{"-verbose", "no"}, false},
		{"0", []string{"-verbose", "0"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type s struct {
				Verbose bool `flag:"verbose"`
			}
			var v s
			if err := FromFlags(context.Background(), tt.args, &v); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertEqual(t, "Verbose", v.Verbose, tt.want)
		})
	}
}

func TestFromFlags_StringSliceSingleElement(t *testing.T) {
	args := []string{"-tags", "only_one"}

	var s flagBasicStruct
	if err := FromFlags(context.Background(), args, &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertSliceEqual(t, "Tags", s.Tags, []string{"only_one"})
}
