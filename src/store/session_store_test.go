package store

import (
	"testing"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSessionStore(t *testing.T) {
	t.Run("create and get session", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		require.NoError(t, TestDb.SessionStore.Create(&model.Session{
			ID:           id,
			Email:        "nguyenvana@gmail.com",
			RefreshToken: "abcd",
			UserAgent:    "Chrome",
			ClientIP:     "192.168.1.1",
			ExpiresAt:    time.Now(),
		}))

		inserted, err := TestDb.SessionStore.GetSession(id)
		require.NoError(t, err)
		require.NotNil(t, inserted)

		require.Equal(t, "nguyenvana@gmail.com", inserted.Email)
		require.Equal(t, "abcd", inserted.RefreshToken)
		require.Equal(t, "Chrome", inserted.UserAgent)
		require.Equal(t, "192.168.1.1", inserted.ClientIP)
	})
}
