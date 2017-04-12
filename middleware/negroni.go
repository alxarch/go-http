package middleware

import "net/http"

type NegroniHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

func Negroni(n NegroniHandler) Middleware {
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		fnext, ok := next.(http.HandlerFunc)
		if !ok {
			fnext = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})

		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n.ServeHTTP(w, r, fnext)
		})
	})
}
