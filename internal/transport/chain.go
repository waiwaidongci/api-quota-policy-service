package transport

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, ms ...Middleware) http.Handler {
	for i := 0; i < len(ms); i++ {
		h = ms[i](h)
	}
	return h
}
