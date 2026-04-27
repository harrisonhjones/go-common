package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointerAndValue(t *testing.T) {
	assert.Equal(t, 42, Value(Pointer(42)))
	assert.Equal(t, 0, Value((*int)(nil)))

	assert.Equal(t, 3.14, Value(Pointer(3.14)))
	assert.Equal(t, (float32)(0.0), Value((*float32)(nil)))

	assert.Equal(t, "foo", Value(Pointer("foo")))
	assert.Equal(t, "", Value((*string)(nil)))
}
