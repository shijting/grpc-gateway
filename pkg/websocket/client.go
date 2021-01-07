package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/showiot/camera/inits/logger"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var upgrader = &websocket.Upgrader{
	//ReadBufferSize:  10240,
	//WriteBufferSize: 10240,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type send struct {
	To          string
	Msg         []byte
	MsgType     int
	ContentType int
}

var clients = make(map[string]*client)

type client struct {
	// The websocket connection.
	conn *websocket.Conn
	// 消息
	sendData chan *send
	// uuid
	id string
	sync.Mutex
}

// 接收消息
type ReceiveMsg struct {
	FromUser    string      `json:"from_user"`
	ToUser      string      `json:"to_user"`
	Content     interface{} `json:"content"`
	ContentType int         `json:"content_type"`
}

func (c *client) writePump(w http.ResponseWriter, r *http.Request) {
	fields := logrus.Fields{
		"func":   "websocket",
		"remote": r.RemoteAddr,
		"uuid":   c.id,
	}
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()
	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case data, ok := <-c.sendData:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// send chan is closed
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			toClient, ok := clients[data.To]
			if !ok {
				logger.GetLogger().WithFields(fields).WithField("to uuid", toClient.id).Error("the client not exist")
				// todo
				continue
			}
			if err := toClient.conn.WriteMessage(websocket.TextMessage, data.Msg); err !=nil {
				logger.GetLogger().WithFields(fields).WithField("to uuid", toClient.id).WithError(err).Error("write message failed")
			}

		}
	}
}
func (c *client) readPump(w http.ResponseWriter, r *http.Request) {
	fields := logrus.Fields{
		"func":   "websocket",
		"remote": r.RemoteAddr,
		"uuid":   c.id,
	}
	defer c.close()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msgByte, err := c.conn.ReadMessage()
		if err != nil {
			//if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
			//	log.Println("client is closed...")
			//	return
			//}
			return
		}
		var receive ReceiveMsg
		if err = json.Unmarshal(msgByte, &receive); err != nil {
			log.Println("unmarshal error:", err)
			logger.GetLogger().WithFields(fields).WithError(err).Error("invalid packet")
			// TODO 返回错误信息
			return
		}
		// 业务处理
		switch receive.ContentType {
		case 1:
		default:
		}
	}
}
func (c *client) close() {
	c.Lock()
	c.conn.Close()
	delete(clients, c.id)
	close(c.sendData)
	c.Unlock()
}
// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	//var mu sync.Mutex
	// get uuid

	id := ""
	fields := logrus.Fields{
		"func":   "websocket",
		"remote": r.RemoteAddr,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		logger.GetLogger().WithFields(fields).WithError(err).Error("upgrade websocket failed")
		return
	}

	client := &client{conn: conn, sendData: make(chan *send, 100), id: id}
	clients[id] = client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump(w, r)
	go client.readPump(w, r)
}
