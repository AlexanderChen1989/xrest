package main

import (
	"fmt"
	"net/http"

	"github.com/AlexanderChen1989/xrest"
	"github.com/AlexanderChen1989/xrest/plug/hello"
	"golang.org/x/net/context"
)

func main() {
	p := xrest.NewPipeline()

	p.Plug(hello.NewHelloHandler("Alex1"))
	p.Plug(hello.NewHelloHandler("Alex2"))
	p.Plug(hello.NewHelloHandler("Alex3"))

	p.HandleFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Println("Done!")
	})

	p.Handler().ServeHTTP(nil, nil)
}
