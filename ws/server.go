package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	sync.RWMutex
	Routes map[string]HandleFunc
	Addr   string

	authentication Authentication

	ConnToUserMp map[*websocket.Conn]string
	UserToConnMp map[string]*websocket.Conn

	Upgrader websocket.Upgrader
}

func NewServer(addr string) *Server {
	return &Server{
		Routes:         make(map[string]HandleFunc),
		Addr:           addr,
		authentication: Authentication{},
		ConnToUserMp:   make(map[*websocket.Conn]string),
		UserToConnMp:   make(map[string]*websocket.Conn),
		Upgrader:       websocket.Upgrader{},
	}
}

func (s *Server) AddConn(conn *websocket.Conn, r *http.Request) {
	uid := s.authentication.GetUid(r)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.ConnToUserMp[conn] = uid
	s.UserToConnMp[uid] = conn
}

func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.UserToConnMp[uid]
}

func (s *Server) GetUsers(conns ...*websocket.Conn) []string {

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

func (s *Server) Close(conn *websocket.Conn) {
	conn.Close()

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.ConnToUserMp[conn]
	delete(s.ConnToUserMp, conn)
	delete(s.UserToConnMp, uid)
}

func (s *Server) SendByUserIds(msg interface{}, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}

	connList := make([]*websocket.Conn, 0, len(userIds))

	for _, id := range userIds {
		connList = append(connList, s.GetConn(id))
	}

	return s.Send(msg, connList...)
}

func (s *Server) Send(msg interface{}, conns ...*websocket.Conn) error {
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

	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade err:%v", err)
		return
	}

	if !s.authentication.Auth(w, r) {
		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("鉴权不通过")))
		return
	}

	//记录链接
	s.AddConn(conn, r)

	go s.HandlerConn(conn)
}

func (s *Server) HandlerConn(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.Close(conn)
			return
		}

		var msg Message

		if err = json.Unmarshal(message, &msg); err != nil {
			s.Close(conn)
			return
		}

		if handle, ok := s.Routes[msg.Method]; ok {
			handle(s, conn, &msg)
		} else {
			err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在的method:%v,请检查", msg.Method)))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

	}
}

func (s *Server) Start() {

	http.HandleFunc("/ws", s.ServerWs)

	fmt.Println("start ws")

	err := http.ListenAndServe(s.Addr, nil)
	if err != nil {
		fmt.Printf("http.ListenAndServe failed: %v", err.Error())
		return
	}
}

func (s *Server) Stop() {
	fmt.Println("stop ws")
}
