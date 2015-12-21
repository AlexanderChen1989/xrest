package xrest

import (
	"net/http"

	"golang.org/x/net/context"
)

type Pipeline struct {
	handler Handler
	plugs   []Plugger
}

func (p *Pipeline) HTTPHandler() http.Handler {
	var h Handler = p.handler
	for i := len(p.plugs) - 1; i >= 0; i-- {
		h = p.plugs[i].Plug(h)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		h.ServeHTTP(ctx, w, r)
	})
}

func (p *Pipeline) Handler() Handler {
	var h Handler = p.handler
	for i := len(p.plugs) - 1; i >= 0; i-- {
		h = p.plugs[i].Plug(h)
	}
	return h
}

func (p *Pipeline) Plug(plugs ...Plugger) *Pipeline {
	p.plugs = append(p.plugs, plugs...)
	return p
}

func (p *Pipeline) Plugs() []Plugger {
	return p.plugs
}

func (p *Pipeline) SetHandler(h Handler) *Pipeline {
	p.handler = h
	return p
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		handler: HandlerFunc(func(context.Context, http.ResponseWriter, *http.Request) {}),
	}
}
