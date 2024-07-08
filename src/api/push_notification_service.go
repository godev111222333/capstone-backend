package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var _ INotificationPushService = (*NotificationPushService)(nil)

type INotificationPushService interface {
}

type PushMessage struct {
	To    []string    `json:"to"`
	Body  string      `json:"body"`
	Title string      `json:"title,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type NotificationPushService struct {
	Url string
}

func NewNotificationPushService() *NotificationPushService {
	return &NotificationPushService{Url: "https://exp.host/--/api/v2/push/send"}
}

func (s *NotificationPushService) Push(m *PushMessage) error {
	bz, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("NotificationPushService: Push %v\n", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.Url, bytes.NewReader(bz))
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if code := resp.StatusCode; code < 200 || code > 299 {
		err := fmt.Errorf("HTTP error. Status code %d", code)
		fmt.Println(err)
		return err
	}

	return nil
}
