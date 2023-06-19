package xsync

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	p := NewPoolR(func(t *int) {
		*t = 1
	}, func(t *int) {
		*t = 2
	})

	v := p.Get()
	require.Equal(t, 1, v.Value)
	v.Release()

	v = p.Get()
	require.Equal(t, 2, v.Value)
	v.Release()

	v = p.Get()
	v = p.Get()
}
