package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) sendMsgToAllJoiners(convID int, content, sender string) {
	joiners, ok := s.chatRooms.Load(convID)
	if !ok {
		fmt.Println("something has wrong")
		return
	}

	jrs, _ := joiners.([]*websocket.Conn)

	fmt.Printf("len room: %d\n", len(jrs))
	for _, conn := range jrs {
		if err := conn.WriteJSON(Message{
			MsgType: MessageTypeTexting,
			Content: content,
			Sender:  sender,
		}); err != nil {
			fmt.Printf("send msg to all joiners err %v\n", err)
		}
	}
}

func (s *Server) removeConnFromRoom(convID int, conn *websocket.Conn) {
	joiners, ok := s.chatRooms.Load(convID)
	if !ok {
		return
	}

	jrs, _ := joiners.([]*websocket.Conn)
	newSubs := make([]*websocket.Conn, 0)
	for _, joiner := range jrs {
		if joiner != conn {
			newSubs = append(newSubs, joiner)
		}
	}

	s.chatRooms.Store(convID, newSubs)
}

func (s *Server) joinConversation(
	convID int,
	conn *websocket.Conn,
) {
	joiners, exist := s.chatRooms.LoadOrStore(convID, []*websocket.Conn{conn})
	if exist {
		jrs, _ := joiners.([]*websocket.Conn)
		jrs = append(jrs, conn)
		s.chatRooms.Store(convID, jrs)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		s.removeConnFromRoom(convID, conn)
		return nil
	})

	s.sendMsgToAllJoiners(convID, "New comer has joined", "system")
}

func sendError(conn *websocket.Conn, err error) error {
	return conn.WriteJSON(Message{
		MsgType: MessageTypeError,
		Content: err.Error(),
	})
}

type (
	MessageType       string
	SystemMessageType string
)

const (
	MessageTypeUserJoin               MessageType = "USER_JOIN"
	MessageTypeAdminJoin              MessageType = "ADMIN_JOIN"
	MessageTypeTexting                MessageType = "TEXTING"
	MessageTypeError                  MessageType = "ERROR"
	MessageTypeSystemResponseUserJoin MessageType = "SYSTEM_USER_JOIN_RESPONSE"
)

type Message struct {
	MsgType        MessageType `json:"msg_type"`
	AccessToken    string      `json:"access_token,omitempty"`
	Content        string      `json:"content,omitempty"`
	ConversationID int         `json:"conversation_id,omitempty"`
	Sender         string      `json:"sender"`
}

func (s *Server) decodeBearerAccessToken(authorize string) (*token.Payload, error) {
	fields := strings.Fields(authorize)
	if len(fields) < 2 {
		err := errors.New("invalid authorization format")
		return nil, err
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		return nil, err
	}

	accessToken := fields[1]
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (s *Server) HandleChat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("upgrade to ws error %v\n", err)
		responseCustomErr(c, ErrCodeUnableUpgradeWebsocket, err)
		return
	}

	go func() {
	loop:
		for {
			msg := Message{}
			if err := conn.ReadJSON(&msg); err != nil {
				break loop
			}

			switch msg.MsgType {
			case MessageTypeAdminJoin:
				if !s.handleAdminJoinMsg(conn, msg) {
					break loop
				}
				break
			case MessageTypeUserJoin:
				_, ok := s.handleUserJoinMsg(conn, msg)
				if !ok {
					break loop
				}

				break
			case MessageTypeTexting:
				sender, ok := s.handleTextingMsg(conn, msg)
				if !ok {
					break loop
				}

				s.adminNewConversationQueue <- ConversationMsg{ConversationID: msg.ConversationID, Sender: sender}
				break
			default:
				fmt.Println("invalid message_type. stop the chat")
				break loop
			}
		}
	}()
}

func (s *Server) checkConnection(convID int, conn *websocket.Conn) {
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()
loop:
	for {
		select {
		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				s.removeConnFromRoom(convID, conn)
				break loop
			}
		}
	}
}

func (s *Server) handleAdminJoinMsg(conn *websocket.Conn, msg Message) bool {
	authPayload, err := s.decodeBearerAccessToken(msg.AccessToken)
	if err != nil {
		fmt.Println(err)
		_ = sendError(conn, err)
		return false
	}
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		_ = sendError(conn, err)
		return false
	}

	if acct.Role.RoleName != model.RoleNameAdmin || acct.Status != model.AccountStatusActive {
		_ = sendError(conn, errors.New("invalid admin or account is inactive"))
		return false
	}

	s.joinConversation(msg.ConversationID, conn)
	go s.checkConnection(msg.ConversationID, conn)
	return true
}

func (s *Server) handleUserJoinMsg(conn *websocket.Conn, msg Message) (int, bool) {
	authPayload, err := s.decodeBearerAccessToken(msg.AccessToken)
	if err != nil {
		fmt.Println(err)
		_ = sendError(conn, err)
		return -1, false
	}
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		_ = sendError(conn, err)
		return -1, false
	}

	if acct.Status != model.AccountStatusActive {
		_ = sendError(conn, errors.New("account is inactive"))
		return -1, false
	}

	conv, err := s.store.ConversationStore.GetByAccID(acct.ID)
	if err != nil {
		_ = sendError(conn, err)
		return -1, false
	}

	if conv == nil {
		conv = &model.Conversation{
			AccountID: acct.ID,
			Status:    model.ConversationStatusActive,
		}

		if err := s.store.ConversationStore.Create(conv); err != nil {
			_ = sendError(conn, err)
			return -1, false
		}
	}

	if err := conn.WriteJSON(Message{
		MsgType:        MessageTypeSystemResponseUserJoin,
		ConversationID: conv.ID,
	}); err != nil {
		fmt.Printf("write JSON %v\v", err)
		_ = sendError(conn, err)
		return -1, false
	}

	s.joinConversation(conv.ID, conn)
	go s.checkConnection(msg.ConversationID, conn)
	return conv.ID, true
}

func (s *Server) handleTextingMsg(conn *websocket.Conn, msg Message) (string, bool) {
	authPayload, err := s.decodeBearerAccessToken(msg.AccessToken)
	if err != nil {
		fmt.Println(err)
		_ = sendError(conn, err)
		return "", false
	}
	acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
	if err != nil {
		_ = sendError(conn, err)
		return "", false
	}

	if authPayload.Role == model.RoleNameAdmin {
		go func() {
			conv, err := s.store.ConversationStore.GetByID(msg.ConversationID)
			if err != nil {
				return
			}

			if conv != nil {
				phone, expoToken := conv.Account.PhoneNumber, s.getExpoToken(conv.Account.PhoneNumber)
				_ = s.notificationPushService.Push(conv.AccountID, s.notificationPushService.NewChatMsg(expoToken, phone))
			}
		}()
	}

	s.sendMsgToAllJoiners(msg.ConversationID, msg.Content, acct.Role.RoleName)
	if err := s.store.MessageStore.Create(&model.Message{
		ConversationID: msg.ConversationID,
		Sender:         acct.ID,
		Content:        msg.Content,
	}); err != nil {
		_ = sendError(conn, err)
		return "", false
	}

	return authPayload.Role, true
}
