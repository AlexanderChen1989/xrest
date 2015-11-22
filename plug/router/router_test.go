package router

import (
	"fmt"
	"testing"
)

func TestRouter(t *testing.T) {
	router := NewRouter()

	sr := router.SubRouter(nil, "/api")
	auth := sr.SubRouter("/auth")
	auth.Get("/files", nil)
	noauth := sr.SubRouter("/noauth")
	noauth.Post("/login", nil)

	for _, route := range router.Routes() {
		fmt.Println(route.Method, route.Path)
	}
}
