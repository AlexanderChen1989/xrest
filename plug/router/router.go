package router

import (
	"net/http"

	"github.com/AlexanderChen1989/xrest"
	"golang.org/x/net/context"
)

type Router struct {
	NotFound  xrest.HandleFunc
	methodMap map[string]map[string]xrest.Handler
}

func (router *Router) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if mm := router.methodMap[r.Method]; mm != nil {
		if handler := mm[r.URL.Path]; handler != nil {
			handler.ServeHTTP(ctx, w, r)
			return
		}
	}

	router.NotFound(ctx, w, r)
}

func (router *Router) Plug(h xrest.Handler) xrest.Handler {
	return router
}

func notFound(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found.", http.StatusNotFound)
}

const GET = "GET"
const POST = "POST"
const PUT = "PUT"
const PATCH = "PATCH"

func NewRouter() *Router {
	mm := make(map[string]map[string]xrest.Handler)

	for _, m := range []string{GET, POST, PUT, PATCH} {
		mm[m] = make(map[string]xrest.Handler)
	}

	router := &Router{
		NotFound:  notFound,
		methodMap: mm,
	}
	return router
}

func (router *Router) handle(method string, path string, h xrest.Handler) {
	router.methodMap[method][path] = h
}

func (router *Router) Get(path string, h xrest.Handler) {
	router.handle(GET, path, h)
}

func (router *Router) Post(path string, h xrest.Handler) {
	router.handle(POST, path, h)
}

func (router *Router) Put(path string, h xrest.Handler) {
	router.handle(PUT, path, h)
}

func (router *Router) Patch(path string, h xrest.Handler) {
	router.handle(PATCH, path, h)
}
