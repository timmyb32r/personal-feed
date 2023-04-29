package model

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComplexKey(t *testing.T) {
	key1 := NewComplexKey("lvl1!blablabla")
	require.Equal(t, 0, key1.Depth())
	key2 := key1.MakeSubkey("lvl2")
	require.Equal(t, 1, key2.Depth())
	require.Equal(t, `lvl1%21blablabla!lvl2`, key2.FullKey())
	key2Restored, err := ParseComplexKey(key2.FullKey())
	require.NoError(t, err)
	require.Equal(t, `lvl1%21blablabla!lvl2`, key2Restored.FullKey())
}
