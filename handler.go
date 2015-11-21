package xrest

import (
	"net/http"

	"golang.org/x/net/context"
)

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

type Plugger interface {
	Plug(Handler) Handler
}

type HandleFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h HandleFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h(ctx, w, r)
}
