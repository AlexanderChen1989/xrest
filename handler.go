package xrest

import (
	"net/http"

	"golang.org/x/net/context"
)

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h HandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h(ctx, w, r)
}

func (h HandlerFunc) Plug(_ Handler) Handler {
	return h
}
