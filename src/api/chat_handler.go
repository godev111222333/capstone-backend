package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
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
		fmt.Println("empty room. an actor joining")
		s.chatRooms.Store(convID, []*websocket.Conn{conn})

		s.sendMsgToAllJoiners(convID, "New comer has joined")
		return
	}

	if jrs, ok := joiners.([]*websocket.Conn); ok {
		fmt.Println("exist >= 1 actor in chat room. an actor joining")
		jrs = append(jrs, conn)
		s.sendMsgToAllJoiners(convID, "New comer has joined")
		s.chatRooms.Store(convID, jrs)
	}
}

func (s *Server) adminJoin(conn *websocket.Conn) {
	s.chatRooms.Range(func(key, value any) bool {
		if convID, ok := key.(int); ok {
			s.joinConversation(convID, conn)
		}

		return true
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

			fmt.Printf("receiving %s msg: %s\n", msg.MsgType, msg.Content)

			switch msg.MsgType {
			case MessageTypeAdminJoin:
				s.adminJoin(conn)
				break
			case MessageTypeUserJoin:
				// TODO: replace acctID with bearer access token
				acctID, err := strconv.Atoi(msg.AccessToken)
				if err != nil {
					fmt.Println(err)
					break loop
				}
				conv := &model.Conversation{
					AccountID: acctID,
					Status:    model.ConversationStatusActive,
				}
				if err := s.store.ConversationStore.Create(conv); err != nil {
					break loop
				}
				if err := conn.WriteJSON(Message{
					MsgType:        MessageTypeSystemResponseUserJoin,
					ConversationID: conv.ID,
				}); err != nil {
					fmt.Println(err)
					break loop
				}

				s.joinConversation(conv.ID, conn)
				break
			case MessageTypeTexting:
				s.sendMsgToAllJoiners(msg.ConversationID, msg.Content)
				s.store.MessageStore.Create(&model.Message{
					ConversationID: msg.ConversationID,
					Sender:         0,
					Content:        msg.Content,
				})
				break
			default:
				fmt.Println("invalid message_type. stop the chat")
				break loop
			}
		}
	}()
}
