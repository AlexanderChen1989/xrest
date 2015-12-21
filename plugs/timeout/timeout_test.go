package timeout

import (
	"net/http"
	"testing"
	"time"

	"github.com/AlexanderChen1989/xrest"
	"golang.org/x/net/context"
)

func TestTimeout(t *testing.T) {
	to := newTimeout(1 * time.Second)

	pipe := xrest.NewPipeline()
	pipe.Plug(to)
	pipe.Plug(xrest.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
		case <-time.After(to.Duration + time.Second):
			t.Error("Shoud timeout\n")
		}
	}))
	pipe.HTTPHandler().ServeHTTP(nil, nil)

	to = newTimeout(0)
	pipe = xrest.NewPipeline()
	pipe.Plug(to)
	pipe.Plug(xrest.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
			t.Error("Shoud not timeout\n")
		case <-time.After(time.Second):
		}
	}))
	pipe.HTTPHandler().ServeHTTP(nil, nil)
}
