package mapstruct

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AlexanderChen1989/xrest"
	"github.com/AlexanderChen1989/xrest/plugs"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/context"
)

// New create new plugger
func New(onError plugs.OnErrorFn) xrest.Plugger {
	p := &plug{
		onError: onError,
	}
	if p.onError == nil {
		p.onError = plugs.DefaultOnErrorFn
	}

	return p
}

type plug struct {
	next    xrest.Handler
	onError plugs.OnErrorFn
}

var ctxKey uint8

func (p *plug) Plug(h xrest.Handler) xrest.Handler {
	p.next = h
	return p
}

func (p *plug) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	payload := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil && p.onError != nil {
		p.onError(ctx, w, r, err)
		return
	}
	p.next.ServeHTTP(context.WithValue(ctx, &ctxKey, payload), w, r)
}

// ErrPlugNotPlugged mapstruct plug is not plugged to pipeline
var ErrPlugNotPlugged = errors.New("mapstruct plug is not plugged to pipeline")

// Decode payload to struct
func Decode(ctx context.Context, val interface{}) error {
	payload, ok := ctx.Value(&ctxKey).(map[string]interface{})
	if !ok {
		return ErrPlugNotPlugged
	}
	mapstructure.Decode(payload, val)
	return nil
}
