package middleware

import (
	"net/http"
	"sync/atomic"
	"time"
)

func OnTick(mw Middleware, d time.Duration) Middleware {
	if d <= 0 {
		return mw
	}
	tick := time.NewTicker(d)
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		total := uint64(0)
		go func() {
			for _ = range tick.C {
				atomic.StoreUint64(&total, 0)
			}
		}()
		wrapped := mw.Wrap(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := atomic.AddUint64(&total, 1)
			if t == 1 {
				wrapped.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})
}
