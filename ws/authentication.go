package ws

import (
	"fmt"
	"im/global"
	"im/middlewares"
	"net/http"
)

type authentication interface {
	Auth(w http.ResponseWriter, r *http.Request) (*middlewares.CustomClaims, bool)
}

type Authentication struct {
}

func (*Authentication) Auth(w http.ResponseWriter, r *http.Request) (*middlewares.CustomClaims, bool) {
	token := r.Header.Get("token")

	fmt.Println(token)

	claims, err := middlewares.Auth(token)

	if err != nil {
		global.Logger["err"].Errorf("Token verification failed, err: %v", err.Error())
		return nil, false
	}

	return claims, true
}
