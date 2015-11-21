package xrest

import (
	"net/http"

	"golang.org/x/net/context"
)

type Pipeline struct {
	handle HandleFunc
	plugs  []Plugger
}

func (p *Pipeline) Handler() http.Handler {
	var h Handler = p.handle
	for i := len(p.plugs) - 1; i >= 0; i-- {
		h = p.plugs[i].Plug(h)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		h.ServeHTTP(ctx, w, r)
	})
}

func (p *Pipeline) Plug(plug Plugger) {
	p.plugs = append(p.plugs, plug)

}

func (p *Pipeline) HandleFunc(fn func(ctx context.Context, w http.ResponseWriter, r *http.Request)) {
	p.handle = HandleFunc(fn)
}

func emptyHandleFunc(context.Context, http.ResponseWriter, *http.Request) {}

func NewPipeline() *Pipeline {
	return &Pipeline{
		handle: emptyHandleFunc,
	}
}
