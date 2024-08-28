package ws

import (
	"github.com/gorilla/websocket"
	"im/global"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	idleMu sync.Mutex
	*websocket.Conn
	userId            int    //当前连接属于哪个用户
	roomId            string //当前用户在哪个房间
	s                 *Server
	idle              time.Time
	maxConnectionIdle time.Duration
	done              chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request, userId int) *Conn {
	c, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		global.Logger["err"].Errorf("upgrade err:%v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		userId:            userId,
		idle:              time.Now(),
		maxConnectionIdle: time.Duration(global.ServerConfig.WsInfo.MaxConnectionIdle) * time.Second,
		done:              make(chan struct{}),
	}

	go conn.keepalive()

	return conn
}

// 长连接检测机制
func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer idleTimer.Stop()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle

			global.Logger["info"].Infof("idle %v, maxIdle %v \n", c.idle, c.maxConnectionIdle)
			if idle.IsZero() { // The connection is non-idle.
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			global.Logger["info"].Infof("val %v \n", val)
			c.idleMu.Unlock()
			if val <= 0 {
				// The connection has been idle for a duration of keepalive.MaxConnectionIdle or more.
				// Gracefully close the connection.
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			global.Logger["info"].Infof("客户端结束连接")
			return
		}
	}
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()
	c.idle = time.Time{}
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	err := c.Conn.WriteMessage(messageType, data)
	// 当写操作完成后当前连接就会进入空闲状态，并记录空闲的时间
	c.idle = time.Now()
	return err
}

func (c *Conn) GetUserId() int {
	return c.userId
}

func (c *Conn) SetUserId(userId int) {
	c.userId = userId
	return
}

func (c *Conn) GetRoomId() string {
	return c.roomId
}

func (c *Conn) SetRoomId(roomId string) {
	c.roomId = roomId
	return
}

func (c *Conn) Close() error {
	close(c.done)
	return c.Conn.Close()
}
