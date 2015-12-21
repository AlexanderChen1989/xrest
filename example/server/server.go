package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/xrest"
	"github.com/AlexanderChen1989/xrest/plugs/body"
	"github.com/AlexanderChen1989/xrest/plugs/close"
	"github.com/AlexanderChen1989/xrest/plugs/limit"
	"github.com/AlexanderChen1989/xrest/plugs/router"
	"github.com/AlexanderChen1989/xrest/plugs/static"
	"github.com/AlexanderChen1989/xrest/utils"
)

func helloRoute(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	params, _ := router.FetchParams(ctx)
	utils.DumpJSON(
		w,
		map[string]interface{}{
			"status": "success",
			"msg":    "Hello, " + params.ByName("name") + "!",
		},
	)
}

func slowRoute(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	utils.DumpJSON(w, `Im slow!`)
}

func newRouter() *router.Router {
	r := router.New()
	r.Get("/api/hello/:name", xrest.HandlerFunc(helloRoute))
	r.Post("/api/hello/:name", xrest.HandlerFunc(helloRoute))
	r.Get("/api/slow", xrest.HandlerFunc(slowRoute))
	return r
}

func main() {
	p := xrest.NewPipeline()

	p.Plug(limit.New(1, time.Second))
	p.Plug(close.New(func(r *http.Request) {
		fmt.Printf("CLOSED: %#v\n", r.RemoteAddr)
	}))
	p.Plug(
		static.New(
			static.Dir("./static"),
			static.Prefix("public"),
		),
	)
	p.Plug(body.New(func(r *http.Request, err error) {
		fmt.Println("Error: ", err)
	}))
	p.Plug(newRouter())

	http.ListenAndServe(":8080", p.HTTPHandler())
}
