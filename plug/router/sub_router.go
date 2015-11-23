package router

import "github.com/AlexanderChen1989/xrest"

type SubRouter struct {
	prefix   string
	pipeline *xrest.Pipeline
	father   *Router
}

func (sr *SubRouter) SubRouter(prefix string) *SubRouter {
	return sr.father.subRouter(sr, prefix)
}

func (sr *SubRouter) handle(method string, path string, h xrest.Handler) {
	sr.father.handle(sr, method, path, h)
}

func (sr *SubRouter) Plug(plug ...xrest.Plugger) *SubRouter {
	sr.father.plug(sr, plug...)
	return sr
}

func (sr *SubRouter) Get(path string, h xrest.Handler) {
	sr.father.Get(path, h)
}

func (sr *SubRouter) Post(path string, h xrest.Handler) {
	sr.father.Post(path, h)
}

func (sr *SubRouter) Put(path string, h xrest.Handler) {
	sr.father.Put(path, h)
}

func (sr *SubRouter) Patch(path string, h xrest.Handler) {
	sr.father.Patch(path, h)
}

func (sr *SubRouter) Options(path string, h xrest.Handler) {
	sr.father.Options(path, h)
}

func (sr *SubRouter) Delete(path string, h xrest.Handler) {
	sr.father.Delete(path, h)
}
