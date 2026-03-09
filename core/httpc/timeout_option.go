package httpc

import (
	"context"
	"net/http"
	"time"
)

type TimeoutOption struct {
	Timeout time.Duration
}

func (t *TimeoutOption) Next(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), t.Timeout)
		r = r.WithContext(ctx)
		defer cancel()
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusRequestTimeout)
			return
		default:
			f(w, r)
		}
	}
}
