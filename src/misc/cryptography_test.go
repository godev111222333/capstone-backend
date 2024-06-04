package misc

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashVerifier(t *testing.T) {
	t.Parallel()

	t.Run("hash and compare hash", func(t *testing.T) {
		t.Parallel()

		h := NewHashVerifier()
		originPass := "admin"
		hashedPass, err := h.Hash(originPass)
		b64, err := base64.StdEncoding.DecodeString(hashedPass)
		require.NoError(t, err)
		require.NoError(t, err)
		require.NoError(t, bcrypt.CompareHashAndPassword(b64, []byte(originPass)))
	})
}
