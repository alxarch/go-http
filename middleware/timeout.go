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
		return bytes.NewBuffer(make([]byte, 0, DefaultBufferSize))
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
			defer func() {
				putBuffer(rec.Body)
			}()
			go func() {
				defer close(done)
				next.ServeHTTP(rec, r.WithContext(ctx))
			}()
			select {
			case <-ctx.Done():
				switch err := ctx.Err(); err {
				case context.DeadlineExceeded:
					timeoutHandler.ServeHTTP(w, r)
				default:
					// Should never be the case
					panic(err)
				}
			case <-done:
				dest := w.Header()
				for k, v := range rec.HeaderMap {
					dest[k] = v
				}
				w.WriteHeader(rec.Code)
				rec.Body.WriteTo(w)
			}
		})
	})

}
