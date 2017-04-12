package middleware

import (
	"mime"
	"net/http"
	"strings"
)

type Middleware interface {
	Wrap(http.Handler) http.Handler
}

type MiddlewareFunc func(http.Handler) http.Handler

func (mf MiddlewareFunc) Wrap(next http.Handler) http.Handler {
	return mf(next)
}

type MiddlewareProvider interface {
	Middleware() Middleware
}

type Middlewares []Middleware

func (ms Middlewares) Wrap(next http.Handler) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
		m := ms[i]
		next = m.Wrap(next)
	}
	return next
}

var NopHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func Compose(mws ...Middleware) Middleware {
	return Middlewares(mws)
}

func RequireMethods(methods ...string) Middleware {
	allowed := make(map[string]bool)
	for _, m := range methods {
		allowed[strings.ToUpper(m)] = true
	}
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allowed[r.Method] {
				next.ServeHTTP(w, r)
				return
			}
			code := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)
		})
	})
}

func RequireMimeTypes(mimes ...string) Middleware {
	allowed := make(map[string]bool)
	for _, m := range mimes {
		allowed[m] = true
	}
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mimetype, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if allowed[mimetype] {
				next.ServeHTTP(w, r)
				return
			}
			code := http.StatusUnsupportedMediaType
			http.Error(w, http.StatusText(code), code)
		})
	})
}
