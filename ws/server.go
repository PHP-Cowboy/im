package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"im/global"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	sync.RWMutex
	Routes map[string]HandleFunc
	Addr   string

	authentication Authentication

	ConnToUserMp map[*Conn]string
	UserToConnMp map[string]*Conn

	Upgrader websocket.Upgrader
}

func NewServer(addr string) *Server {
	return &Server{
		Routes:         make(map[string]HandleFunc),
		Addr:           addr,
		authentication: Authentication{},
		ConnToUserMp:   make(map[*Conn]string),
		UserToConnMp:   make(map[string]*Conn),
		Upgrader:       websocket.Upgrader{},
	}
}

func (s *Server) AddConn(conn *Conn, r *http.Request) {
	uid := s.authentication.GetUid(r)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 原有已经存在了连接
	if c := s.UserToConnMp[uid]; c != nil {
		c.Close()
	}

	s.ConnToUserMp[conn] = uid
	s.UserToConnMp[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.UserToConnMp[uid]
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.ConnToUserMp))
		for _, uid := range s.ConnToUserMp {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.ConnToUserMp[conn])
		}
	}

	return res
}

// 关闭链接
func (s *Server) Close(conn *Conn) {

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.ConnToUserMp[conn]
	if uid == "" {
		// 已经关闭了连接
		return
	}

	delete(s.ConnToUserMp, conn)
	delete(s.UserToConnMp, uid)

	conn.Close()
}

func (s *Server) SendByUserIds(msg interface{}, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}

	connList := make([]*Conn, 0, len(userIds))

	for _, id := range userIds {
		connList = append(connList, s.GetConn(id))
	}

	return s.Send(msg, connList...)
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil

}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.Routes[r.Method] = r.Handler
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("server handle recover failed, err:%v", err)
		}
	}()

	conn := NewConn(s, w, r)

	if !s.authentication.Auth(w, r) {

		if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("鉴权不通过"))); err != nil {
			global.Logger["err"].Error(err.Error())
		}

		return
	}

	//记录链接
	s.AddConn(conn, r)

	go s.handlerConn(conn)
}

func (s *Server) handlerConn(conn *Conn) {
	// 记录连接
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// 关闭并删除连接
			s.Close(conn)
			return
		}

		// 请求信息
		var message Message
		json.Unmarshal(msg, &message)

		// 依据请求消息类型分类处理
		switch message.FrameType {
		case FramePing:
			// ping：回复
			s.Send(&Message{FrameType: FramePing}, conn)
		case FrameData:
			// 处理
			if handler, ok := s.Routes[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				s.Send(&Message{
					FrameType: FrameData,
					Data:      fmt.Sprintf("不存在请求方法 %v 请仔细检查", message.Method),
				}, conn)
			}
		}
	}
}

func (s *Server) Start() {

	http.HandleFunc("/", s.ServerWs)

	fmt.Println("start ws")

	err := http.ListenAndServe(s.Addr, nil)
	if err != nil {
		global.Logger["err"].Errorf("http.ListenAndServe failed: %v", err.Error())
		return
	}
}

func (s *Server) Stop() {
	fmt.Println("stop ws")
}
