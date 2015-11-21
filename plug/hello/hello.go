package hello

import (
	"fmt"
	"net/http"

	"github.com/AlexanderChen1989/xrest"

	"golang.org/x/net/context"
)

type HelloHandler struct {
	next xrest.Handler
	name string
}

func (hello *HelloHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello", hello.name)
	hello.next.ServeHTTP(ctx, w, r)
}

func (hello *HelloHandler) Plug(h xrest.Handler) xrest.Handler {
	hello.next = h
	return hello
}

func NewHelloHandler(name string) *HelloHandler {
	return &HelloHandler{name: name}
}
