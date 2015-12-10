package timeout

import (
	"net/http"
	"time"

	"github.com/AlexanderChen1989/xrest"
	"golang.org/x/net/context"
)

// timeout plug
// if duration <= 0, no timeout
type timeout struct {
	next     xrest.Handler
	Duration time.Duration
}

// New create a new timeout plug
func New(d time.Duration) xrest.Plugger {
	return newTimeout(d)
}

func newTimeout(d time.Duration) *timeout {
	return &timeout{Duration: d}
}

// Plug implements xrest.Plugger interface
func (to *timeout) Plug(h xrest.Handler) xrest.Handler {
	to.next = h
	return to
}

// ServeHTTP implements xrest.Handler interface
func (to *timeout) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if to.Duration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, to.Duration)
		defer cancel()
	}
	to.next.ServeHTTP(ctx, w, r)
}
