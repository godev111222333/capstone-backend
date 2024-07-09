package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/gorilla/websocket"
)

const AdminSubsKey = -1

type NotificationMsg struct {
	Title string      `json:"title,omitempty"`
	Body  string      `json:"body,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func (s *Server) startAdminSub() {
	go func() {
		for {
			select {
			case msg := <-s.adminNotificationQueue:
				s.sendAdminNotification(msg)
			}
		}
	}()

	// test function
	go func() {
		testTicker := time.NewTicker(3 * time.Second)
		for {
			select {
			case <-testTicker.C:
				s.adminNotificationQueue <- NotificationMsg{
					Title: "test title",
					Body:  "test body",
					Data:  map[string]string{"aa": "bb"},
				}
			}
		}

	}()
}

func (s *Server) sendAdminNotification(msg NotificationMsg) {
	subs, ok := s.adminNotificationSubs.Load(AdminSubsKey)
	if !ok {
		return
	}

	toSubs, _ := subs.([]*websocket.Conn)
	for _, sub := range toSubs {
		if err := sub.WriteJSON(msg); err != nil {
			s.removeAdminNotificationSub(sub)
		}
	}
}

func (s *Server) checkAdminConnection(conn *websocket.Conn) {
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()
loop:
	for {
		select {
		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				s.removeAdminNotificationSub(conn)
				break loop
			}
		}
	}
}

func (s *Server) removeAdminNotificationSub(conn *websocket.Conn) {
	subs, ok := s.adminNotificationSubs.Load(AdminSubsKey)
	if ok {
		curSubs, _ := subs.([]*websocket.Conn)
		newSubs := make([]*websocket.Conn, 0)
		for _, e := range curSubs {
			if e != conn {
				newSubs = append(newSubs, e)
			}
		}

		s.adminNotificationSubs.Store(AdminSubsKey, newSubs)
	}
}

func (s *Server) adminSubscribeNotification(conn *websocket.Conn) {
	curSubscribers, isLoaded := s.adminNotificationSubs.LoadOrStore(AdminSubsKey, []*websocket.Conn{conn})
	if curSubs, ok := curSubscribers.([]*websocket.Conn); ok && isLoaded {
		curSubscribers = append(curSubs, conn)
		s.adminNotificationSubs.Store(AdminSubsKey, curSubscribers)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		s.removeAdminNotificationSub(conn)
		return nil
	})
}

type AuthMsg struct {
	AccessToken string `json:"access_token"`
}

func (s *Server) HandleAdminSubscribeNotification(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		responseCustomErr(c, ErrCodeUnableUpgradeWebsocket, err)
		return
	}

	auth := AuthMsg{}
	if err := conn.ReadJSON(&auth); err != nil {
		_ = sendError(conn, err)
		return
	}

	authPayload, err := s.decodeBearerAccessToken(auth.AccessToken)
	if err != nil {
		_ = sendError(conn, err)
		return
	}

	if authPayload.Role != model.RoleNameAdmin {
		_ = sendError(conn, errors.New("invalid role"))
		return
	}

	s.adminSubscribeNotification(conn)
	go s.checkAdminConnection(conn)
}
