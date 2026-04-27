package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMust(t *testing.T) {
	require.Equal(t, 1, Must(func() (int, error) {
		return 1, nil
	}()))
	require.Equal(t, "a", Must(func() (string, error) {
		return "a", nil
	}()))
	require.Equal(t, 1.0, Must(func() (float64, error) {
		return 1.0, nil
	}()))
	require.PanicsWithError(t, "boom", func() {
		Must(func() (any, error) {
			return nil, fmt.Errorf("boom")
		}())
	})
}
