package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

type bodyStore struct {
	pool *sync.Pool
}

var ctxBodyKey uint8

type readCloser struct {
	io.ReadCloser
	bs  *bodyStore
	buf *bytes.Buffer
}

func NewBodyStore() *bodyStore {
	return &bodyStore{
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
	}
}

var defaultBodyStore = NewBodyStore()

var DecodeJSON = defaultBodyStore.DecodeJSON

func (rc *readCloser) Close() error {
	rc.bs.pool.Put(rc.buf)

	return rc.ReadCloser.Close()
}

func (bs *bodyStore) DecodeJSON(ctx context.Context, r *http.Request, v interface{}) (context.Context, error) {
	// fetch a buf from pool
	data, ok := ctx.Value(&ctxBodyKey).([]byte)
	if !ok {
		buf := bs.pool.Get().(*bytes.Buffer)
		buf.Reset()
		if _, err := io.Copy(buf, r.Body); err != nil {
			bs.pool.Put(buf)
			return ctx, err
		}
		data = buf.Bytes()
		ctx = context.WithValue(ctx, &ctxBodyKey, data)
		// reconstruct http.Request.Body
		rc := &readCloser{
			ReadCloser: r.Body,
			bs:         bs,
			buf:        buf,
		}
		r.Body = rc
	}

	return ctx, json.Unmarshal(data, v)
}
