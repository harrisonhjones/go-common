package common

// Value dereferences the pointer and returns the underlying value.
// If the pointer is nil, the zero value of T is returned.
func Value[T any](v *T) T {
	if v == nil {
		var t T
		return t
	}
	return *v
}

// Pointer returns a pointer to the value.
func Pointer[T any](v T) *T {
	return &v
}
