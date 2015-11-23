package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/xrest"
	"github.com/AlexanderChen1989/xrest/plugs/hello"
)

func TestRouter(t *testing.T) {

	router := NewRouter()

	sr := router.SubRouter("/api")

	auth := sr.SubRouter("/auth")
	noauth := sr.SubRouter("/noauth")

	auth.Plug(hello.NewHelloHandler("Hello1"))
	auth.Plug(hello.NewHelloHandler("Hello2"))
	noauth.Plug(hello.NewHelloHandler("World1"))
	noauth.Plug(hello.NewHelloHandler("World2"))

	auth.Get("/files", xrest.HandleFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth files Done!")
	}))
	noauth.Post("/login", xrest.HandleFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Println("Notauth login Done!")
	}))

	authReq, _ := http.NewRequest("GET", "/api/auth/files", nil)
	noauthReq, _ := http.NewRequest("POST", "/api/noauth/login", nil)

	pipe := xrest.NewPipeline().Plug(router)

	w := httptest.NewRecorder()

	pipe.HTTPHandler().ServeHTTP(w, authReq)
	pipe.HTTPHandler().ServeHTTP(w, noauthReq)
}
