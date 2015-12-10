package body

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/AlexanderChen1989/xrest"

	"golang.org/x/net/context"
)

type bodyPlug struct {
	pool *sync.Pool
	next xrest.Handler
}

type readCloser struct {
	io.ReadCloser
	bp  *bodyPlug
	buf *bytes.Buffer
}

func newBody() *bodyPlug {
	return &bodyPlug{
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
	}
}

func New() xrest.Plugger {
	return newBody()
}

// Close return buf to pool
func (rc *readCloser) Close() error {
	rc.bp.pool.Put(rc.buf)

	return rc.ReadCloser.Close()
}

// ErrBodyNotPlugged body plug not plugged
var ErrBodyNotPlugged = errors.New("Body not plugged.")

var ctxBodyKey uint8

// DecodeJSON decode json to interface from body
func DecodeJSON(ctx context.Context, v interface{}) error {
	data, ok := ctx.Value(&ctxBodyKey).([]byte)

	if !ok {
		return ErrBodyNotPlugged
	}

	return json.Unmarshal(data, v)
}

func FetchBody(ctx context.Context) ([]byte, bool) {
	body, ok := ctx.Value(&ctxBodyKey).([]byte)
	return body, ok
}

func (bp *bodyPlug) Plug(h xrest.Handler) xrest.Handler {
	bp.next = h
	return bp
}

func (bp *bodyPlug) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if _, ok := FetchBody(ctx); !ok {
		buf := bp.pool.Get().(*bytes.Buffer)
		buf.Reset()
		if _, err := io.Copy(buf, r.Body); err != nil {
			bp.pool.Put(buf)
		}
		ctx = context.WithValue(ctx, &ctxBodyKey, buf.Bytes())
		// reconstruct http.Request.Body
		rc := &readCloser{
			ReadCloser: r.Body,
			bp:         bp,
			buf:        buf,
		}
		r.Body = rc
	}

	bp.next.ServeHTTP(ctx, w, r)
}
