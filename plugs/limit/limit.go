package limit

import (
	"net/http"
	"time"

	"github.com/AlexanderChen1989/xrest"
	"github.com/AlexanderChen1989/xrest/utils"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
	"github.com/didip/tollbooth/errors"
	"golang.org/x/net/context"
)

// limiter cancel context when timeout
type limiter struct {
	limiter *config.Limiter

	ErrHandleFn func(http.ResponseWriter, *errors.HTTPError)

	next xrest.Handler
}

func errHandleFn(w http.ResponseWriter, err *errors.HTTPError) {
	utils.JSON.Code(err.StatusCode).Dump(w, map[string]interface{}{
		"status": "error",
		"msg":    err.Message,
	})
}

// New create a new request rate limiter plug, max requests in ttl time duration
func New(max int64, ttl time.Duration) xrest.Plugger {
	return &limiter{
		limiter:     tollbooth.NewLimiter(max, ttl),
		ErrHandleFn: errHandleFn,
	}
}

// NewLimiter create a new request rate limiter with conf and error handle function
func NewLimiter(conf *config.Limiter, handlefn func(http.ResponseWriter, *errors.HTTPError)) xrest.Plugger {
	return &limiter{
		limiter:     conf,
		ErrHandleFn: handlefn,
	}
}

// ServeHTTP implement Handler.ServeHTTP
func (limiter *limiter) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := tollbooth.LimitByRequest(limiter.limiter, r)
	if err != nil {
		limiter.ErrHandleFn(w, err)
		return
	}

	// There's no rate-limit error, serve the next handler.
	limiter.next.ServeHTTP(ctx, w, r)
}

// Plug implement Plugger.Plug
func (limiter *limiter) Plug(h xrest.Handler) xrest.Handler {
	limiter.next = h
	return limiter
}
