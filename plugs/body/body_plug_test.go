package body

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/AlexanderChen1989/xrest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type testPayload struct {
	Name string
	Age  int
}

type testPlug struct {
	t      *testing.T
	next   xrest.Handler
	origin testPayload
}

func newTestPlug(t *testing.T, origin testPayload) *testPlug {
	return &testPlug{t: t, origin: origin}
}

func (tp *testPlug) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var payload testPayload
	err := DecodeJSON(ctx, &payload)
	assert.Nil(tp.t, err)
	assert.True(tp.t, reflect.DeepEqual(tp.origin, payload))
	tp.next.ServeHTTP(ctx, w, r)
}

func (tp *testPlug) Plug(h xrest.Handler) xrest.Handler {
	tp.next = h
	return tp
}

func TestJSONDecodeIntegration(t *testing.T) {
	pipe := xrest.NewPipeline()
	origin := testPayload{
		Name: "alex",
		Age:  27,
	}
	pipe.Plug(
		New(nil),
		newTestPlug(t, origin),
		newTestPlug(t, origin),
		newTestPlug(t, origin),
		newTestPlug(t, origin),
	)
	var buf bytes.Buffer
	assert.Nil(t, json.NewEncoder(&buf).Encode(&origin))
	r, _ := http.NewRequest("GET", "/api/test", &buf)
	pipe.HTTPHandler().ServeHTTP(nil, r)
}
