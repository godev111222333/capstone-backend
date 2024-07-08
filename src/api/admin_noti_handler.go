package api

import (
	"github.com/gorilla/websocket"
)

type NotificationMsg struct {
	Title string
	Body  string
	Data  interface{}
}

func (s *Server) adminSubscribeNotification(conn *websocket.Conn) {
	curSubscribers, isLoaded := s.adminNotificationSubs.LoadOrStore(-1, []*websocket.Conn{conn})
	if curSubs, ok := curSubscribers.([]*websocket.Conn); ok && isLoaded {
		curSubscribers = append(curSubs, conn)
		s.adminNotificationSubs.Store(-1, curSubscribers)
	}

	//conn.SetCloseHandler(func(code int, text string) error {
	//	if ok {
	//		newSubs := make([]*websocket.Conn, 0)
	//		for _, subscriber := range curSubs {
	//			if subscriber != conn {
	//				newSubs = append(newSubs, subscriber)
	//			}
	//		}
	//		s.adminNotificationSubs.Store(-1, newSubs)
	//	}
	//	return nil
	//})
}

//func (s *Server) HandleAdminSubscribeNotification(c *gin.Context) {
//	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		fmt.Println(err)
//		responseCustomErr(c, ErrCodeUnableUpgradeWebsocket, err)
//		return
//	}
//
//	go func() {
//	loop:
//		for {
//		}
//	}()
//}
