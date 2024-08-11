package handler

import (
	"im/ws"
)

func RegisterHandlers(s *ws.Server) {
	s.AddRoutes([]ws.Route{
		{
			Method:  "user.online",
			Handler: Online(),
		},
	})
}
