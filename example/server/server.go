package main

import (
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
	utils.WriteJSON(
		w,
		http.StatusOK,
		map[string]interface{}{
			"status": "success",
			"msg":    "Hello, " + params.ByName("name") + "!",
		},
	)
}

func newRouter() *router.Router {
	r := router.New()
	r.Get("/api/hello/:name", xrest.HandleFunc(helloRoute))
	return r
}

func main() {
	p := xrest.NewPipeline()

	p.Plug(limit.New(1, time.Second))
	p.Plug(close.New())
	p.Plug(
		static.New(
			static.Dir("./static"),
			static.Prefix("public"),
		),
	)
	p.Plug(body.New())
	p.Plug(newRouter())

	http.ListenAndServe(":8080", p.HTTPHandler())
}
