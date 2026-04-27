package common

// Must returns v if err is nil, and panics with err if err is non-nil.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
