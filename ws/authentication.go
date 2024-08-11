package ws

import (
	"net/http"
	"strconv"
	"time"
)

type authentication interface {
	Auth(w http.ResponseWriter, r *http.Request) bool
	GetUid(r *http.Request) string
}

type Authentication struct {
}

func (*Authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func (*Authentication) GetUid(r *http.Request) string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}
