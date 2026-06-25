package ratelimiters

import (
	"net/http"
	"strings"
)

type Ratelimter interface {
	Allow(r *http.Request) bool
}

type GetIp func(r *http.Request) string

func DefaultGetIp(splitPort bool) GetIp {
	return func(r *http.Request) string {
		if splitPort {
			return strings.Split(r.RemoteAddr, ":")[0]
		} else {
			return r.RemoteAddr
		}
	}
}