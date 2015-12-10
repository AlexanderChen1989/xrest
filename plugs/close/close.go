package close

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/xrest"
)

// Closer client connection close plug
type closer struct {
	next xrest.Handler
}

// New create closer plug
func New() xrest.Plugger {
	return &closer{}
}

// Plug implements xrest.Plugger interface
func (c *closer) Plug(h xrest.Handler) xrest.Handler {
	c.next = h
	return c
}

// ServeHTTP implements xrest.Handler interface
func (c *closer) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Cancel the context if the client closes the connection
	if cn, ok := w.(http.CloseNotifier); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		go func() {
			<-cn.CloseNotify()
			cancel()
		}()
	}

	c.next.ServeHTTP(ctx, w, r)
}
