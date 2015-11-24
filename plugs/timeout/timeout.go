package timeout

import (
	"net/http"
	"time"

	"github.com/AlexanderChen1989/xrest"
	"golang.org/x/net/context"
)

// Timeout timeout plug
// if duration <= 0, no timeout
type Timeout struct {
	next     xrest.Handler
	Duration time.Duration
}

// New create a new timeout plug
func New(d time.Duration) *Timeout {
	return &Timeout{Duration: d}
}

// Plug implements Plugger interface
func (to *Timeout) Plug(h xrest.Handler) xrest.Handler {
	to.next = h
	return to
}

// ServeHTTP implements xrest.Handler interface
func (to *Timeout) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if to.Duration > 0 {
		ctx, _ = context.WithTimeout(ctx, to.Duration)
	}
	to.next.ServeHTTP(ctx, w, r)
}
