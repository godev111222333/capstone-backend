package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNotificationPushService(t *testing.T) {
	service := NewNotificationPushService()
	require.NoError(t, service.Push(&PushMessage{
		To:    []string{},
		Body:  "Body noti",
		Title: "Title noti",
		Data: map[string]interface{}{
			"url": "https://google.com",
		},
	}))
}
