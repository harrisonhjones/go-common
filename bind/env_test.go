package bind

import (
	"context"
	"fmt"
	"os"
	"testing"
)

// --- test structs ---

type basicStruct struct {
	EnvString      string   `env:"TEST_STRING"`
	EnvBoolTrue    bool     `env:"TEST_BOOL_TRUE"`
	EnvBoolY       bool     `env:"TEST_BOOL_Y"`
	EnvBoolYes     bool     `env:"TEST_BOOL_YES"`
	EnvBoolFalse   bool     `env:"TEST_BOOL_FALSE"`
	EnvInt         int      `env:"TEST_INT"`
	EnvInt8        int8     `env:"TEST_INT8"`
	EnvInt16       int16    `env:"TEST_INT16"`
	EnvInt32       int32    `env:"TEST_INT32"`
	EnvInt64       int64    `env:"TEST_INT64"`
	EnvUint        uint     `env:"TEST_UINT"`
	EnvUint8       uint8    `env:"TEST_UINT8"`
	EnvUint16      uint16   `env:"TEST_UINT16"`
	EnvUint32      uint32   `env:"TEST_UINT32"`
	EnvUint64      uint64   `env:"TEST_UINT64"`
	EnvFloat32     float32  `env:"TEST_FLOAT32"`
	EnvFloat64     float64  `env:"TEST_FLOAT64"`
	EnvStringSlice []string `env:"TEST_STRING_SLICE"`
}

type defaultsStruct struct {
	EnvString      string   `env:"TEST_UNSET_STRING" envDefault:"default_string"`
	EnvBool        bool     `env:"TEST_UNSET_BOOL" envDefault:"true"`
	EnvInt         int      `env:"TEST_UNSET_INT" envDefault:"-42"`
	EnvInt8        int8     `env:"TEST_UNSET_INT8" envDefault:"12"`
	EnvInt16       int16    `env:"TEST_UNSET_INT16" envDefault:"-13"`
	EnvInt32       int32    `env:"TEST_UNSET_INT32" envDefault:"14"`
	EnvInt64       int64    `env:"TEST_UNSET_INT64" envDefault:"-15"`
	EnvUint        uint     `env:"TEST_UNSET_UINT" envDefault:"1"`
	EnvUint8       uint8    `env:"TEST_UNSET_UINT8" envDefault:"2"`
	EnvUint16      uint16   `env:"TEST_UNSET_UINT16" envDefault:"3"`
	EnvUint32      uint32   `env:"TEST_UNSET_UINT32" envDefault:"4"`
	EnvUint64      uint64   `env:"TEST_UNSET_UINT64" envDefault:"5"`
	EnvFloat32     float32  `env:"TEST_UNSET_FLOAT32" envDefault:"1.2"`
	EnvFloat64     float64  `env:"TEST_UNSET_FLOAT64" envDefault:"3.4"`
	EnvStringSlice []string `env:"TEST_UNSET_SLICE" envDefault:"foo,bar,baz"`
}

type noTagStruct struct {
	Name string
	Age  int
}

type mixedTagStruct struct {
	Tagged   string `env:"TEST_MIXED_TAGGED"`
	Untagged string
}

type normalizerStruct struct {
	Name       string `env:"TEST_NORM_NAME"`
	normalized bool
}

func (n *normalizerStruct) Normalize(_ context.Context) {
	n.normalized = true
	n.Name = "normalized:" + n.Name
}

type validatorStruct struct {
	Name string `env:"TEST_VALID_NAME"`
}

func (v *validatorStruct) Validate(_ context.Context) error {
	if v.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type normAndValidStruct struct {
	Name       string `env:"TEST_NV_NAME"`
	normalized bool
}

func (nv *normAndValidStruct) Normalize(_ context.Context) {
	nv.normalized = true
}

func (nv *normAndValidStruct) Validate(_ context.Context) error {
	if nv.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

type envOverrideDefaultStruct struct {
	Val string `env:"TEST_OVERRIDE" envDefault:"fallback"`
}

type badIntStruct struct {
	Val int `env:"TEST_BAD_INT"`
}

type badUintStruct struct {
	Val uint `env:"TEST_BAD_UINT"`
}

type badFloatStruct struct {
	Val float64 `env:"TEST_BAD_FLOAT"`
}

type unsupportedFieldStruct struct {
	Val complex128 `env:"TEST_COMPLEX"`
}

type unsupportedSliceStruct struct {
	Val []int `env:"TEST_INT_SLICE"`
}

// --- helpers ---

func setEnvs(t *testing.T, envs map[string]string) {
	t.Helper()
	for k, v := range envs {
		t.Setenv(k, v)
	}
}

// --- tests ---

func TestFromEnv_NonPointer(t *testing.T) {
	var s basicStruct
	err := FromEnv(context.Background(), s)
	if err == nil {
		t.Fatal("expected error for non-pointer, got nil")
	}
	if err.Error() != "v must be a pointer" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromEnv_PointerToNonStruct(t *testing.T) {
	s := "hello"
	err := FromEnv(context.Background(), &s)
	if err == nil {
		t.Fatal("expected error for pointer to non-struct, got nil")
	}
	if err.Error() != "v must be a pointer to a struct" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromEnv_BasicTypes(t *testing.T) {
	setEnvs(t, map[string]string{
		"TEST_STRING":       "hello",
		"TEST_BOOL_TRUE":    "true",
		"TEST_BOOL_Y":       "y",
		"TEST_BOOL_YES":     "yes",
		"TEST_BOOL_FALSE":   "false",
		"TEST_INT":          "-100",
		"TEST_INT8":         "127",
		"TEST_INT16":        "-32000",
		"TEST_INT32":        "2147483647",
		"TEST_INT64":        "-9223372036854775808",
		"TEST_UINT":         "100",
		"TEST_UINT8":        "255",
		"TEST_UINT16":       "65535",
		"TEST_UINT32":       "4294967295",
		"TEST_UINT64":       "18446744073709551615",
		"TEST_FLOAT32":      "3.14",
		"TEST_FLOAT64":      "2.718281828",
		"TEST_STRING_SLICE": "a,b,c",
	})

	var s basicStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "EnvString", s.EnvString, "hello")
	assertEqual(t, "EnvBoolTrue", s.EnvBoolTrue, true)
	assertEqual(t, "EnvBoolY", s.EnvBoolY, true)
	assertEqual(t, "EnvBoolYes", s.EnvBoolYes, true)
	assertEqual(t, "EnvBoolFalse", s.EnvBoolFalse, false)
	assertEqual(t, "EnvInt", s.EnvInt, -100)
	assertEqual(t, "EnvInt8", s.EnvInt8, int8(127))
	assertEqual(t, "EnvInt16", s.EnvInt16, int16(-32000))
	assertEqual(t, "EnvInt32", s.EnvInt32, int32(2147483647))
	assertEqual(t, "EnvInt64", s.EnvInt64, int64(-9223372036854775808))
	assertEqual(t, "EnvUint", s.EnvUint, uint(100))
	assertEqual(t, "EnvUint8", s.EnvUint8, uint8(255))
	assertEqual(t, "EnvUint16", s.EnvUint16, uint16(65535))
	assertEqual(t, "EnvUint32", s.EnvUint32, uint32(4294967295))
	assertEqual(t, "EnvUint64", s.EnvUint64, uint64(18446744073709551615))
	assertFloat(t, "EnvFloat32", float64(s.EnvFloat32), 3.14, 0.01)
	assertFloat(t, "EnvFloat64", s.EnvFloat64, 2.718281828, 0.000001)
	assertSliceEqual(t, "EnvStringSlice", s.EnvStringSlice, []string{"a", "b", "c"})
}

func TestFromEnv_Defaults(t *testing.T) {
	// Ensure none of the TEST_UNSET_* vars are set (t.Setenv not called)
	var s defaultsStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "EnvString", s.EnvString, "default_string")
	assertEqual(t, "EnvBool", s.EnvBool, true)
	assertEqual(t, "EnvInt", s.EnvInt, -42)
	assertEqual(t, "EnvInt8", s.EnvInt8, int8(12))
	assertEqual(t, "EnvInt16", s.EnvInt16, int16(-13))
	assertEqual(t, "EnvInt32", s.EnvInt32, int32(14))
	assertEqual(t, "EnvInt64", s.EnvInt64, int64(-15))
	assertEqual(t, "EnvUint", s.EnvUint, uint(1))
	assertEqual(t, "EnvUint8", s.EnvUint8, uint8(2))
	assertEqual(t, "EnvUint16", s.EnvUint16, uint16(3))
	assertEqual(t, "EnvUint32", s.EnvUint32, uint32(4))
	assertEqual(t, "EnvUint64", s.EnvUint64, uint64(5))
	assertFloat(t, "EnvFloat32", float64(s.EnvFloat32), 1.2, 0.01)
	assertFloat(t, "EnvFloat64", s.EnvFloat64, 3.4, 0.01)
	assertSliceEqual(t, "EnvStringSlice", s.EnvStringSlice, []string{"foo", "bar", "baz"})
}

func TestFromEnv_EnvOverridesDefault(t *testing.T) {
	t.Setenv("TEST_OVERRIDE", "from_env")

	var s envOverrideDefaultStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", s.Val, "from_env")
}

func TestFromEnv_DefaultUsedWhenEnvUnset(t *testing.T) {
	// Don't set TEST_OVERRIDE
	var s envOverrideDefaultStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", s.Val, "fallback")
}

func TestFromEnv_NoTags(t *testing.T) {
	var s noTagStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Fields without tags should remain zero-valued
	assertEqual(t, "Name", s.Name, "")
	assertEqual(t, "Age", s.Age, 0)
}

func TestFromEnv_MixedTags(t *testing.T) {
	t.Setenv("TEST_MIXED_TAGGED", "tagged_value")

	var s mixedTagStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Tagged", s.Tagged, "tagged_value")
	assertEqual(t, "Untagged", s.Untagged, "")
}

func TestFromEnv_EmptyStruct(t *testing.T) {
	type empty struct{}
	var s empty
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromEnv_BoolVariants(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"true", "true", true},
		{"True", "True", true},
		{"TRUE", "TRUE", true},
		{"y", "y", true},
		{"Y", "Y", true},
		{"yes", "yes", true},
		{"Yes", "Yes", true},
		{"YES", "YES", true},
		{"false", "false", false},
		{"0", "0", false},
		{"no", "no", false},
		{"empty", "", false},
		{"random", "random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				t.Setenv("TEST_BOOL_VAR", tt.value)
			} else {
				// For empty string, set it explicitly
				os.Setenv("TEST_BOOL_VAR", tt.value)
				t.Cleanup(func() { os.Unsetenv("TEST_BOOL_VAR") })
			}

			type boolStruct struct {
				Val bool `env:"TEST_BOOL_VAR"`
			}
			var s boolStruct
			if err := FromEnv(context.Background(), &s); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertEqual(t, "Val", s.Val, tt.want)
		})
	}
}

func TestFromEnv_InvalidInt(t *testing.T) {
	t.Setenv("TEST_BAD_INT", "not_a_number")
	var s badIntStruct
	err := FromEnv(context.Background(), &s)
	if err == nil {
		t.Fatal("expected error for invalid int, got nil")
	}
}

func TestFromEnv_InvalidUint(t *testing.T) {
	t.Setenv("TEST_BAD_UINT", "-1")
	var s badUintStruct
	err := FromEnv(context.Background(), &s)
	if err == nil {
		t.Fatal("expected error for invalid uint, got nil")
	}
}

func TestFromEnv_InvalidFloat(t *testing.T) {
	t.Setenv("TEST_BAD_FLOAT", "not_a_float")
	var s badFloatStruct
	err := FromEnv(context.Background(), &s)
	if err == nil {
		t.Fatal("expected error for invalid float, got nil")
	}
}

func TestFromEnv_UnsupportedFieldType(t *testing.T) {
	t.Setenv("TEST_COMPLEX", "1+2i")
	var s unsupportedFieldStruct
	err := FromEnv(context.Background(), &s)
	if err == nil {
		t.Fatal("expected error for unsupported field type, got nil")
	}
}

func TestFromEnv_UnsupportedSliceType(t *testing.T) {
	t.Setenv("TEST_INT_SLICE", "1,2,3")
	var s unsupportedSliceStruct
	err := FromEnv(context.Background(), &s)
	if err == nil {
		t.Fatal("expected error for unsupported slice element type, got nil")
	}
}

func TestFromEnv_Normalizer(t *testing.T) {
	t.Setenv("TEST_NORM_NAME", "alice")

	var s normalizerStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.normalized {
		t.Fatal("expected Normalize to be called")
	}
	assertEqual(t, "Name", s.Name, "normalized:alice")
}

func TestFromEnv_ValidatorPass(t *testing.T) {
	t.Setenv("TEST_VALID_NAME", "bob")

	var s validatorStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Name", s.Name, "bob")
}

func TestFromEnv_ValidatorFail(t *testing.T) {
	// Don't set TEST_VALID_NAME so Name stays empty
	var s validatorStruct
	err := FromEnv(context.Background(), &s)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestFromEnv_NormalizerAndValidator(t *testing.T) {
	t.Setenv("TEST_NV_NAME", "charlie")

	var s normAndValidStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.normalized {
		t.Fatal("expected Normalize to be called")
	}
	assertEqual(t, "Name", s.Name, "charlie")
}

func TestFromEnv_NormalizerAndValidatorFail(t *testing.T) {
	// Don't set TEST_NV_NAME
	var s normAndValidStruct
	err := FromEnv(context.Background(), &s)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	// Normalize should still have been called before Validate
	if !s.normalized {
		t.Fatal("expected Normalize to be called before Validate")
	}
}

func TestFromEnv_StringSliceSingleElement(t *testing.T) {
	t.Setenv("TEST_STRING_SLICE", "only_one")

	var s basicStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertSliceEqual(t, "EnvStringSlice", s.EnvStringSlice, []string{"only_one"})
}

func TestFromEnv_StringSliceEmpty(t *testing.T) {
	// When env var is set to empty string, FromEnv treats it as unset (val == "")
	// and skips the field since there's no envDefault either.
	os.Setenv("TEST_STRING_SLICE", "")
	t.Cleanup(func() { os.Unsetenv("TEST_STRING_SLICE") })

	var s basicStruct
	if err := FromEnv(context.Background(), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.EnvStringSlice != nil {
		t.Errorf("expected nil slice for empty env, got %v", s.EnvStringSlice)
	}
}

func TestFromEnv_NilPointer(t *testing.T) {
	err := FromEnv(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil, got nil")
	}
}

// --- envFromNoOverwrite tests ---

func TestFromEnv_NoOverwrite_StringPreserved(t *testing.T) {
	t.Setenv("TEST_NO_OW_STR", "from_env")

	type s struct {
		Val string `env:"TEST_NO_OW_STR" envFromNoOverwrite:"true"`
	}
	v := s{Val: "already_set"}
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, "already_set")
}

func TestFromEnv_NoOverwrite_StringZeroLoadsEnv(t *testing.T) {
	t.Setenv("TEST_NO_OW_STR", "from_env")

	type s struct {
		Val string `env:"TEST_NO_OW_STR" envFromNoOverwrite:"true"`
	}
	var v s
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, "from_env")
}

func TestFromEnv_NoOverwrite_IntPreserved(t *testing.T) {
	t.Setenv("TEST_NO_OW_INT", "99")

	type s struct {
		Val int `env:"TEST_NO_OW_INT" envFromNoOverwrite:"true"`
	}
	v := s{Val: 42}
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, 42)
}

func TestFromEnv_NoOverwrite_IntZeroLoadsEnv(t *testing.T) {
	t.Setenv("TEST_NO_OW_INT", "99")

	type s struct {
		Val int `env:"TEST_NO_OW_INT" envFromNoOverwrite:"true"`
	}
	var v s
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, 99)
}

func TestFromEnv_NoOverwrite_BoolTruePreserved(t *testing.T) {
	t.Setenv("TEST_NO_OW_BOOL", "false")

	type s struct {
		Val bool `env:"TEST_NO_OW_BOOL" envFromNoOverwrite:"true"`
	}
	v := s{Val: true}
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, true)
}

func TestFromEnv_NoOverwrite_BoolZeroLoadsEnv(t *testing.T) {
	t.Setenv("TEST_NO_OW_BOOL", "true")

	type s struct {
		Val bool `env:"TEST_NO_OW_BOOL" envFromNoOverwrite:"true"`
	}
	var v s
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, true)
}

func TestFromEnv_NoOverwrite_SlicePreserved(t *testing.T) {
	t.Setenv("TEST_NO_OW_SLICE", "x,y,z")

	type s struct {
		Val []string `env:"TEST_NO_OW_SLICE" envFromNoOverwrite:"true"`
	}
	v := s{Val: []string{"existing"}}
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertSliceEqual(t, "Val", v.Val, []string{"existing"})
}

func TestFromEnv_NoOverwrite_SliceNilLoadsEnv(t *testing.T) {
	t.Setenv("TEST_NO_OW_SLICE", "x,y,z")

	type s struct {
		Val []string `env:"TEST_NO_OW_SLICE" envFromNoOverwrite:"true"`
	}
	var v s
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertSliceEqual(t, "Val", v.Val, []string{"x", "y", "z"})
}

func TestFromEnv_NoOverwrite_WithDefault(t *testing.T) {
	// Field has a value, env is unset, default exists — field should be preserved
	type s struct {
		Val string `env:"TEST_NO_OW_DEF_UNSET" envDefault:"fallback" envFromNoOverwrite:"true"`
	}
	v := s{Val: "preset"}
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, "preset")
}

func TestFromEnv_NoOverwrite_ZeroWithDefault(t *testing.T) {
	// Field is zero, env is unset, default exists — default should be used
	type s struct {
		Val string `env:"TEST_NO_OW_DEF_UNSET2" envDefault:"fallback" envFromNoOverwrite:"true"`
	}
	var v s
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, "fallback")
}

func TestFromEnv_NoOverwrite_NotSet_FieldOverwritten(t *testing.T) {
	// Without envFromNoOverwrite, env value should overwrite existing field value
	t.Setenv("TEST_OW_STR", "from_env")

	type s struct {
		Val string `env:"TEST_OW_STR"`
	}
	v := s{Val: "already_set"}
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, "from_env")
}

func TestFromEnv_NoOverwrite_Float64Preserved(t *testing.T) {
	t.Setenv("TEST_NO_OW_F64", "9.99")

	type s struct {
		Val float64 `env:"TEST_NO_OW_F64" envFromNoOverwrite:"true"`
	}
	v := s{Val: 1.23}
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertFloat(t, "Val", v.Val, 1.23, 0.001)
}

func TestFromEnv_NoOverwrite_UintPreserved(t *testing.T) {
	t.Setenv("TEST_NO_OW_UINT", "999")

	type s struct {
		Val uint `env:"TEST_NO_OW_UINT" envFromNoOverwrite:"true"`
	}
	v := s{Val: 7}
	if err := FromEnv(context.Background(), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Val", v.Val, uint(7))
}

// --- assertion helpers ---

func assertEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func assertFloat(t *testing.T, name string, got, want, epsilon float64) {
	t.Helper()
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Errorf("%s: got %v, want %v (epsilon %v)", name, got, want, epsilon)
	}
}

func assertSliceEqual(t *testing.T, name string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: length mismatch: got %d, want %d (%v vs %v)", name, len(got), len(want), got, want)
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("%s[%d]: got %q, want %q", name, i, got[i], want[i])
		}
	}
}
