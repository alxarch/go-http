package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

const DefaultBufferSize = 8192

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, DefaultBufferSize))
	},
}

func getBuffer() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}
func putBuffer(b *bytes.Buffer) {
	b.Reset()
	pool.Put(b)
}

func Timeout(d time.Duration, timeoutHandler http.Handler) Middleware {
	if nil == timeoutHandler {
		code := http.StatusGatewayTimeout
		msg := http.StatusText(code)
		timeoutHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, msg, code)
		})
	}
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			done := make(chan struct{})
			rec := &httptest.ResponseRecorder{
				HeaderMap: http.Header{},
				Body:      getBuffer(),
				Code:      200,
			}
			go func() {
				defer close(done)
				next.ServeHTTP(rec, r.WithContext(ctx))
			}()
			select {
			case <-ctx.Done():
				switch err := ctx.Err(); err {
				case context.DeadlineExceeded:
					timeoutHandler.ServeHTTP(w, r)
				case context.Canceled:
					http.Error(w, err.Error(), http.StatusInternalServerError)
				default:
					panic(err)
					// pass
				}
				return
			case <-done:
				for k, v := range rec.HeaderMap {
					w.Header()[k] = v
				}
				w.WriteHeader(rec.Code)
				rec.Body.WriteTo(w)
				return
			}
		})
	})

}
