package router

import (
	"net/http"
	"testing"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/xrest"
	"github.com/stretchr/testify/assert"
)

type testPlug struct {
	name string
	next xrest.Handler
}

func newTestPlug(name string) *testPlug {
	return &testPlug{name: name}
}

func (plug *testPlug) Plug(h xrest.Handler) xrest.Handler {
	plug.next = h
	return plug
}

var ctxTestPlug uint8

func (plug *testPlug) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vals, ok := ctx.Value(&ctxTestPlug).([]string)
	if !ok {
		vals = []string{}
	}
	vals = append(vals, plug.name)
	ctx = context.WithValue(ctx, &ctxTestPlug, vals)
	plug.next.ServeHTTP(ctx, w, r)
}

func TestRouter(t *testing.T) {

	router := New()

	sr := router.SubRouter("/api")

	auth := sr.SubRouter("/auth")
	noauth := sr.SubRouter("/noauth")

	authVals := []string{"Hello1", "Hello2", "Hello3"}
	for _, val := range authVals {
		auth.Plug(newTestPlug(val))
	}

	noauthVals := []string{"World1", "World2", "World3"}
	for _, val := range noauthVals {
		noauth.Plug(newTestPlug(val))
	}

	auth.Get("/files", xrest.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		assert.EqualValues(t, authVals, ctx.Value(&ctxTestPlug))
	}))
	noauth.Post("/login", xrest.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		assert.EqualValues(t, noauthVals, ctx.Value(&ctxTestPlug))
	}))

	authReq, _ := http.NewRequest("GET", "/api/auth/files", nil)
	noauthReq, _ := http.NewRequest("POST", "/api/noauth/login", nil)

	pipe := xrest.NewPipeline().Plug(router)

	pipe.HTTPHandler().ServeHTTP(nil, authReq)
	pipe.HTTPHandler().ServeHTTP(nil, noauthReq)
}
