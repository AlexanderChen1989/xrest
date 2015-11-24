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

func newBodyPlug() *bodyPlug {
	return &bodyPlug{
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
	}
}

var DefaultPlug = newBodyPlug()

var DecodeJSON = DefaultPlug.DecodeJSON
var Body = DefaultPlug.Body

// Close return buf to pool
func (rc *readCloser) Close() error {
	rc.bp.pool.Put(rc.buf)

	return rc.ReadCloser.Close()
}

var ErrBodyNotPlugged = errors.New("Body not plugged.")

func (bp *bodyPlug) Body(ctx context.Context) ([]byte, error) {
	data, ok := ctx.Value(&ctxBodyKey).([]byte)

	if !ok {
		return nil, ErrBodyNotPlugged
	}

	return data, nil
}

func (bp *bodyPlug) DecodeJSON(ctx context.Context, r *http.Request, v interface{}) error {
	// fetch a buf from pool
	data, ok := ctx.Value(&ctxBodyKey).([]byte)

	if !ok {
		return ErrBodyNotPlugged
	}

	return json.Unmarshal(data, v)
}

var ctxBodyKey uint8

func (bp *bodyPlug) Plug(h xrest.Handler) xrest.Handler {
	bp.next = h
	return bp
}

func FetchBodyFromCtx(ctx context.Context) ([]byte, bool) {
	body, ok := ctx.Value(&ctxBodyKey).([]byte)
	return body, ok
}

func (bp *bodyPlug) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if _, ok := FetchBodyFromCtx(ctx); !ok {
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
