package router

import (
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/xrest"
)

const (
	GET   = "GET"
	POST  = "POST"
	PUT   = "PUT"
	PATCH = "PATCH"
)

type Router struct {
	methodMap map[string]map[string]xrest.Handler
	subs      map[string]*xrest.Pipeline

	NotFound xrest.HandleFunc
}

func NewRouter() *Router {
	mm := map[string]map[string]xrest.Handler{
		GET:   map[string]xrest.Handler{},
		POST:  map[string]xrest.Handler{},
		PUT:   map[string]xrest.Handler{},
		PATCH: map[string]xrest.Handler{},
	}

	subs := map[string]*xrest.Pipeline{}

	return &Router{
		methodMap: mm,
		subs:      subs,
		NotFound: xrest.HandleFunc(func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
			http.Error(w, "Not found.", http.StatusNotFound)
		}),
	}
}

func (r *Router) Plug(_ xrest.Handler) xrest.Handler {
	return r
}

func (r *Router) plug(sub *SubRouter, plug ...xrest.Plugger) {
	r.subs[sub.prefix].Plug(plug...)
}

func (r *Router) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	if m := r.methodMap[req.Method]; m != nil {
		if h := m[req.URL.Path]; h != nil {
			h.ServeHTTP(ctx, w, req)
			return
		}
	}
	r.NotFound(ctx, w, req)
}

type Route struct {
	Method string
	Path   string
}

func (r *Router) Routes() []Route {
	var routes []Route
	for method, m := range r.methodMap {
		for p, _ := range m {
			routes = append(routes, Route{Method: method, Path: p})
		}
	}
	return routes
}

func (r *Router) handle(sub *SubRouter, method string, path string, h xrest.Handler) {
	if sub != nil {
		path = filepath.Join(sub.prefix, path)
		pipe := r.subs[sub.prefix]
		h = pipe.SetHandler(h).Handler()
	}

	r.methodMap[method][path] = h
}

func (r *Router) Get(path string, h xrest.Handler) {
	r.handle(nil, GET, path, h)
}

func (r *Router) Post(path string, h xrest.Handler) {
	r.handle(nil, POST, path, h)
}

func (r *Router) Put(path string, h xrest.Handler) {
	r.handle(nil, PUT, path, h)
}

func (r *Router) Patch(path string, h xrest.Handler) {
	r.handle(nil, PATCH, path, h)
}

func (r *Router) SubRouter(pre *SubRouter, prefix string) *SubRouter {
	prefix = filepath.Join("/", prefix)

	var plugs []xrest.Plugger
	if pre != nil {
		prefix = filepath.Join(pre.prefix, prefix)
		plugs = r.subs[pre.prefix].Plugs()
	}
	sub := &SubRouter{
		prefix: filepath.Join("/", prefix),
		father: r,
	}

	r.subs[sub.prefix] = xrest.NewPipeline().Plug(plugs...)
	return sub
}
