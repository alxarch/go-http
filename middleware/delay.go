package middleware

import (
	"net/http"
	"time"
)

type Delay time.Duration

func (d Delay) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.NewTimer(time.Duration(d))
		defer t.Stop()
		select {
		case <-t.C:
			next.ServeHTTP(w, r)
		}
	})
}
