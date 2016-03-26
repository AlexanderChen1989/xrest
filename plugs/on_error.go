package plugs

import (
	"log"
	"net/http"

	"golang.org/x/net/context"
)

// OnErrorFn handle error in plug
type OnErrorFn func(context.Context, http.ResponseWriter, *http.Request, error)

// DefaultOnErrorFn default function to handler error in plug
func DefaultOnErrorFn(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
}
