package ws

type Route struct {
	Method  string
	Handler HandleFunc
}

type HandleFunc func(s *Server, conn *Conn, msg *Message)
