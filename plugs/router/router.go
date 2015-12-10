package router

import (
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/xrest"
	"github.com/AlexanderChen1989/xrest/plugs/router/tree"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	DELETE  = "DELETE"
)

type Router struct {
	trees map[string]*tree.Node
	subs  map[string]*xrest.Pipeline

	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool

	MethodNotAllowed xrest.Handler
	NotFound         xrest.Handler
}

func New() *Router {
	return &Router{
		trees: map[string]*tree.Node{},
		subs:  map[string]*xrest.Pipeline{},
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
	}
}

func (r *Router) Plug(_ xrest.Handler) xrest.Handler {
	return r
}

func (r *Router) plug(sub *SubRouter, plug ...xrest.Plugger) {
	r.subs[sub.prefix].Plug(plug...)
}

var ctxParamsKey uint8

func FetchParams(ctx context.Context) (ps tree.Params, ok bool) {
	ps, ok = ctx.Value(&ctxParamsKey).(tree.Params)
	return
}

func (r *Router) ServeHTTP(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	if root := r.trees[req.Method]; root != nil {
		path := req.URL.Path

		if handle, ps, tsr := root.GetValue(path); handle != nil {
			// inject params to context
			ctx = context.WithValue(ctx, &ctxParamsKey, ps)
			handle.ServeHTTP(ctx, w, req)
			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.FindCaseInsensitivePath(
					tree.CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	// Handle 405
	if r.HandleMethodNotAllowed {
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == req.Method {
				continue
			}

			handle, _, _ := r.trees[method].GetValue(req.URL.Path)
			if handle != nil {
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed.ServeHTTP(ctx, w, req)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(ctx, w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) handle(sub *SubRouter, method string, path string, h xrest.Handler) {
	if sub != nil {
		path = filepath.Join(sub.prefix, path)
		pipe := r.subs[sub.prefix]
		h = pipe.SetHandler(h).Handler()
	}
	root := r.trees[method]
	if root == nil {
		root = new(tree.Node)
		r.trees[method] = root
	}
	root.AddRoute(path, h)
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

func (r *Router) Options(path string, h xrest.Handler) {
	r.handle(nil, OPTIONS, path, h)
}

func (r *Router) Delete(path string, h xrest.Handler) {
	r.handle(nil, DELETE, path, h)
}

func (r *Router) SubRouter(prefix string) *SubRouter {
	return r.subRouter(nil, prefix)
}

func (r *Router) subRouter(pre *SubRouter, prefix string) *SubRouter {
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
