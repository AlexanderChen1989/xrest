package router

import "github.com/AlexanderChen1989/xrest"

type SubRouter struct {
	prefix   string
	pipeline *xrest.Pipeline
	father   *Router
}

func (sr *SubRouter) SubRouter(prefix string) *SubRouter {
	sub := sr.father.SubRouter(sr, prefix)
	return sub
}

func (sr *SubRouter) handle(method string, path string, h xrest.Handler) {
	sr.father.handle(sr, method, path, h)
}

func (sr *SubRouter) Plug(plug xrest.Plugger) *SubRouter {
	sr.father.plug(sr, plug)
	return sr
}

func (sr *SubRouter) Get(path string, h xrest.Handler) {
	sr.handle(GET, path, h)
}

func (sr *SubRouter) Post(path string, h xrest.Handler) {
	sr.handle(POST, path, h)
}

func (sr *SubRouter) Put(path string, h xrest.Handler) {
	sr.handle(PUT, path, h)
}

func (sr *SubRouter) Patch(path string, h xrest.Handler) {
	sr.handle(PATCH, path, h)
}
