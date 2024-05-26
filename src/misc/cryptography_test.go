package misc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashVerifier(t *testing.T) {
	t.Parallel()

	t.Run("hash and compare hash", func(t *testing.T) {
		t.Parallel()

		h := NewHashVerifier()
		originPass := "4444"
		hashedPass, err := h.Hash(originPass)
		require.NoError(t, err)
		require.Nil(t, h.Compare(hashedPass, originPass))
		require.Error(t, h.Compare(hashedPass, "wrongpassword"))
	})
}
