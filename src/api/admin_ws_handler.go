package api

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/gorilla/websocket"
)

const (
	AdminNotificationSubsKey = -1
	AdminConversationSubsKey = -2
)

const (
	TypeNotification int = iota + 1
	TypeConversation
)

type NotificationMsg struct {
	Title string      `json:"title,omitempty"`
	Body  string      `json:"body,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type ConversationMsg struct {
	ConversationID int `json:"conversation_id"`
}

type AuthMsg struct {
	AccessToken string `json:"access_token"`
}

func (s *Server) HandleAdminSubscribeNotification(c *gin.Context) {
	conn, err := s.initAdminWsConnection(c)
	if err != nil {
		return
	}

	s.adminSubscribeTo(conn, TypeNotification)
	go s.checkAdminConnection(conn, TypeNotification)
}

func (s *Server) HandleAdminSubscribeNewConversation(c *gin.Context) {
	conn, err := s.initAdminWsConnection(c)
	if err != nil {
		return
	}

	s.adminSubscribeTo(conn, TypeConversation)
	go s.checkAdminConnection(conn, TypeConversation)
}

func (s *Server) initAdminWsConnection(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		responseCustomErr(c, ErrCodeUnableUpgradeWebsocket, err)
		return nil, err
	}

	auth := AuthMsg{}
	if err := conn.ReadJSON(&auth); err != nil {
		_ = sendError(conn, err)
		return nil, err
	}

	authPayload, err := s.decodeBearerAccessToken(auth.AccessToken)
	if err != nil {
		_ = sendError(conn, err)
		return nil, err
	}

	if authPayload.Role != model.RoleNameAdmin {
		err := errors.New("invalid role")
		_ = sendError(conn, err)
		return nil, err
	}

	return conn, nil
}

func (s *Server) startAdminSub() {
	go func() {
		for {
			select {
			case msg := <-s.adminNotificationQueue:
				s.sendAdminNotification(msg, TypeNotification)
				break
			case msg := <-s.adminNewConversationQueue:
				s.sendAdminNotification(msg, TypeConversation)
				break
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
					Body:  "this is test msg from server for testing purpose",
					Data:  map[string]string{"aa": "bb"},
				}

				s.adminNewConversationQueue <- ConversationMsg{
					ConversationID: 1,
				}
			}
		}

	}()
}

func (s *Server) sendAdminNotification(msg interface{}, tz int) {
	key := AdminNotificationSubsKey
	if tz == TypeConversation {
		key = AdminConversationSubsKey
	}
	subs, ok := s.adminSubs.Load(key)
	if !ok {
		return
	}

	toSubs, _ := subs.([]*websocket.Conn)
	for _, sub := range toSubs {
		if err := sub.WriteJSON(msg); err != nil {
			s.removeAdminSub(sub, tz)
		}
	}
}

func (s *Server) checkAdminConnection(conn *websocket.Conn, typez int) {
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()
loop:
	for {
		select {
		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				s.removeAdminSub(conn, typez)
				break loop
			}
		}
	}
}

func (s *Server) removeAdminSub(conn *websocket.Conn, tz int) {
	key := AdminNotificationSubsKey
	if tz == TypeConversation {
		key = AdminConversationSubsKey
	}

	subs, ok := s.adminSubs.Load(key)
	if ok {
		curSubs, _ := subs.([]*websocket.Conn)
		newSubs := make([]*websocket.Conn, 0)
		for _, e := range curSubs {
			if e != conn {
				newSubs = append(newSubs, e)
			}
		}

		s.adminSubs.Store(key, newSubs)
	}
}

func (s *Server) adminSubscribeTo(conn *websocket.Conn, tz int) {
	key := AdminConversationSubsKey
	if tz == TypeConversation {
		key = AdminConversationSubsKey
	}

	curSubscribers, isLoaded := s.adminSubs.LoadOrStore(key, []*websocket.Conn{conn})
	if curSubs, ok := curSubscribers.([]*websocket.Conn); ok && isLoaded {
		curSubscribers = append(curSubs, conn)
		s.adminSubs.Store(key, curSubscribers)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		s.removeAdminSub(conn, tz)
		return nil
	})
}
