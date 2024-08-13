package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"im/global"
	"im/middlewares"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	sync.RWMutex
	Routes map[string]HandleFunc
	Addr   string

	authentication authentication

	ConnToUserMp map[*Conn]int
	UserToConnMp map[int]*Conn

	Upgrader websocket.Upgrader
}

func NewServer(addr string) *Server {
	return &Server{
		Routes:         make(map[string]HandleFunc),
		Addr:           addr,
		authentication: &Authentication{},
		ConnToUserMp:   make(map[*Conn]int),
		UserToConnMp:   make(map[int]*Conn),
		Upgrader:       websocket.Upgrader{},
	}
}

func (s *Server) AddConn(conn *Conn, claims *middlewares.CustomClaims) {
	uid := claims.ID

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 原有已经存在了连接
	if c := s.UserToConnMp[uid]; c != nil {
		c.Close()
	}

	s.ConnToUserMp[conn] = uid
	s.UserToConnMp[uid] = conn
}

func (s *Server) GetConn(uid int) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.UserToConnMp[uid]
}

func (s *Server) GetUsers(conns ...*Conn) []int {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []int
	if len(conns) == 0 {
		// 获取全部
		res = make([]int, 0, len(s.ConnToUserMp))
		for _, uid := range s.ConnToUserMp {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]int, 0, len(conns))
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

	uid, ok := s.ConnToUserMp[conn]
	if !ok {
		// 已经关闭了连接
		return
	}

	delete(s.ConnToUserMp, conn)
	delete(s.UserToConnMp, uid)

	conn.Close()
}

func (s *Server) SendByUserIds(msg interface{}, userIds ...int) error {
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
		global.Logger["err"].Errorf("JSON parsing failed, err: %v ", err.Error())
		return err
	}

	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			global.Logger["err"].Errorf("WriteMessage failed, err: %v ", err.Error())
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

	claims, ok := s.authentication.Auth(w, r)

	if !ok {

		if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("鉴权不通过"))); err != nil {
			global.Logger["err"].Errorf(err.Error())
			//conn.Close()
		}

		return
	}

	//记录链接
	s.AddConn(conn, claims)

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
		err = json.Unmarshal(msg, &message)
		if err != nil {
			s.Send(&Message{
				FrameType: FrameData,
				Data:      fmt.Sprintf("消息解析失败"),
			}, conn)

			global.Logger["err"].Errorf("消息解析失败:%v", err.Error())

			return
		}

		// 依据请求消息类型分类处理
		switch message.FrameType {
		case FramePing:
			// ping：回复
			if err = s.Send(&Message{FrameType: FramePing}, conn); err != nil {
				global.Logger["err"].Errorf("ping：回复 failed, err: %v ", err.Error())
				return
			}
		case FrameData:
			// 处理
			if handler, ok := s.Routes[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				s.Send(&Message{
					FrameType: FrameData,
					Data:      fmt.Sprintf("不存在请求方法 %v 请仔细检查", message.Method),
				}, conn)

				global.Logger["err"].Errorf("不存在请求方法 %v 请仔细检查", message.Method)
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
