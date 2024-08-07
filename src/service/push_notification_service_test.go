package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNotificationPushService(t *testing.T) {
	t.Skip()

	service := NewNotificationPushService("", nil)
	require.NoError(t, service.Push(1, &PushMessage{
		To:    []string{},
		Body:  "Body noti",
		Title: "Title noti",
		Data: map[string]interface{}{
			"url": "https://google.com",
		},
	}))
}
