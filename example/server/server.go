package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/xrest"
	"github.com/AlexanderChen1989/xrest/plugs/body"
	"github.com/AlexanderChen1989/xrest/plugs/router"
)

func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	params, _ := router.FetchParams(ctx)
	fmt.Fprintf(w, "Hello, %s!\n", params.ByName("name"))
}

func main() {
	p := xrest.NewPipeline()
	p.Plug(body.Default)
	r := router.New()
	r.Get("/api/hello/:name", xrest.HandleFunc(hello))
	p.Plug(r)

	http.ListenAndServe(":8080", p.HTTPHandler())
}
