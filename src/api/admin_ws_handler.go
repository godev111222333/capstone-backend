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
	s.processWsWithKey(c, AdminNotificationSubsKey)
}

func (s *Server) HandleAdminSubscribeNewConversation(c *gin.Context) {
	s.processWsWithKey(c, AdminConversationSubsKey)
}

func (s *Server) processWsWithKey(c *gin.Context, key int) {
	conn, err := s.initAdminWsConnection(c)
	if err != nil {
		return
	}

	s.adminSubscribeTo(conn, key)
	go s.checkAdminConnection(conn, key)
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
				s.sendMsgToAdmin(msg, AdminNotificationSubsKey)
				break
			case msg := <-s.adminNewConversationQueue:
				s.sendMsgToAdmin(msg, AdminConversationSubsKey)
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
			}
		}
	}()
}

func (s *Server) sendMsgToAdmin(msg interface{}, key int) {
	subs, ok := s.adminSubs.Load(key)
	if !ok {
		return
	}

	toSubs, _ := subs.([]*websocket.Conn)
	for _, sub := range toSubs {
		if err := sub.WriteJSON(msg); err != nil {
			s.removeAdminSub(sub, key)
		}
	}
}

func (s *Server) checkAdminConnection(conn *websocket.Conn, key int) {
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()
loop:
	for {
		select {
		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				s.removeAdminSub(conn, key)
				break loop
			}
		}
	}
}

func (s *Server) removeAdminSub(conn *websocket.Conn, key int) {
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

func (s *Server) adminSubscribeTo(conn *websocket.Conn, key int) {
	curSubscribers, isLoaded := s.adminSubs.LoadOrStore(key, []*websocket.Conn{conn})
	if curSubs, ok := curSubscribers.([]*websocket.Conn); ok && isLoaded {
		curSubscribers = append(curSubs, conn)
		s.adminSubs.Store(key, curSubscribers)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		s.removeAdminSub(conn, key)
		return nil
	})
}
