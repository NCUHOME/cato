package httpc

import (
	"net/http"
	"time"
)

type TimeoutOption struct {
	Timeout time.Duration
}

func (t *TimeoutOption) Next(f http.HandlerFunc) http.HandlerFunc {
	if t == nil || t.Timeout <= 0 {
		return f
	}
	return defaultTimeoutController().Wrap(t.Timeout, f)
}
