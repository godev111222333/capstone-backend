package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/token"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) sendMsgToAllJoiners(convID int, content string) {
	joiners, ok := s.chatRooms.Load(convID)
	if !ok {
		fmt.Println("something has wrong")
		return
	}

	if jrs, ok := joiners.([]*websocket.Conn); ok {
		for _, conn := range jrs {
			if err := conn.WriteJSON(Message{
				MsgType: MessageTypeTexting,
				Content: content,
			}); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (s *Server) joinConversation(
	convID int,
	conn *websocket.Conn,
) {
	joiners, exist := s.chatRooms.Load(convID)
	if !exist {
		s.chatRooms.Store(convID, []*websocket.Conn{conn})

		s.sendMsgToAllJoiners(convID, "New comer has joined")
		return
	}

	if jrs, ok := joiners.([]*websocket.Conn); ok {
		jrs = append(jrs, conn)
		s.sendMsgToAllJoiners(convID, "New comer has joined")
		s.chatRooms.Store(convID, jrs)
	}
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
}

type SystemMessage struct {
	MsgType        SystemMessageType `json:"system_msg_type"`
	ConversationID int               `json:"conversation_id,omitempty"`
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
		fmt.Println(err)
		responseCustomErr(c, ErrCodeUnableUpgradeWebsocket, err)
		return
	}

	go func() {
	loop:
		for {
			msg := Message{}
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Println(err)
				break loop
			}

			switch msg.MsgType {
			case MessageTypeAdminJoin:
				authPayload, err := s.decodeBearerAccessToken(msg.AccessToken)
				if err != nil {
					fmt.Println(err)
					_ = sendError(conn, err)
					break
				}
				acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
				if err != nil {
					fmt.Println(err)
					_ = sendError(conn, err)
					break
				}

				if acct.Role.RoleName != model.RoleNameAdmin || acct.Status != model.AccountStatusActive {
					_ = sendError(conn, errors.New("invalid admin or account is inactive"))
					break
				}

				s.joinConversation(msg.ConversationID, conn)
				break
			case MessageTypeUserJoin:
				authPayload, err := s.decodeBearerAccessToken(msg.AccessToken)
				if err != nil {
					fmt.Println(err)
					_ = sendError(conn, err)
					break
				}
				acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
				if err != nil {
					fmt.Println(err)
					_ = sendError(conn, err)
					break
				}

				if acct.Status != model.AccountStatusActive {
					_ = sendError(conn, errors.New("account is inactive"))
					break
				}

				conv, err := s.store.ConversationStore.GetByAccID(acct.ID)
				if err != nil {
					_ = sendError(conn, err)
					break loop
				}

				if conv == nil {
					conv = &model.Conversation{
						AccountID: acct.ID,
						Status:    model.ConversationStatusActive,
					}

					if err := s.store.ConversationStore.Create(conv); err != nil {
						_ = sendError(conn, err)
						break loop
					}
				}

				if err := conn.WriteJSON(Message{
					MsgType:        MessageTypeSystemResponseUserJoin,
					ConversationID: conv.ID,
				}); err != nil {
					fmt.Println(err)
					_ = sendError(conn, err)
					break loop
				}

				s.joinConversation(conv.ID, conn)
				break
			case MessageTypeTexting:
				authPayload, err := s.decodeBearerAccessToken(msg.AccessToken)
				if err != nil {
					fmt.Println(err)
					_ = sendError(conn, err)
					break
				}
				acct, err := s.store.AccountStore.GetByPhoneNumber(authPayload.PhoneNumber)
				if err != nil {
					fmt.Println(err)
					_ = sendError(conn, err)
					break
				}

				s.sendMsgToAllJoiners(msg.ConversationID, msg.Content)
				if err := s.store.MessageStore.Create(&model.Message{
					ConversationID: msg.ConversationID,
					Sender:         acct.ID,
					Content:        msg.Content,
				}); err != nil {
					_ = sendError(conn, err)
					break loop
				}
				break
			default:
				fmt.Println("invalid message_type. stop the chat")
				break loop
			}
		}
	}()
}
