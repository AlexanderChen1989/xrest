package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlexanderChen1989/xrest"
	"golang.org/x/net/context"
)

func main() {
	p := xrest.NewPipeline()

	p.Plug(xrest.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// fmt.Println("Hello")
	}))

	start := time.Now()
	p.HTTPHandler().ServeHTTP(nil, nil)
	fmt.Println(time.Now().Sub(start))
}
