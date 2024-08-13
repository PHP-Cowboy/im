package ws

import (
	"github.com/golang-jwt/jwt"
	"im/middlewares"
	"math/rand"
	"net/http"
)

type authentication interface {
	Auth(w http.ResponseWriter, r *http.Request) (*middlewares.CustomClaims, bool)
}

type Authentication struct {
}

func (*Authentication) Auth(w http.ResponseWriter, r *http.Request) (*middlewares.CustomClaims, bool) {
	//token := r.Header.Get("token")
	//
	//claims, err := middlewares.Auth(token)
	//
	//if err != nil {
	//	global.Logger["err"].Errorf("Token verification failed, err: %v", err.Error())
	//	return nil, false
	//}

	claims := &middlewares.CustomClaims{
		ID:             rand.Int(),
		Username:       "a",
		StandardClaims: jwt.StandardClaims{},
	}

	return claims, true
}
