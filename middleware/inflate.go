package middleware

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
)

var Inflate = MiddlewareFunc(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header["Content-Encoding"]) == 0 {
			next.ServeHTTP(w, r)
			return
		}
		switch encoding := r.Header.Get("Content-Encoding"); encoding {
		case "gzip":
			if gzr, err := gzip.NewReader(r.Body); err == nil {
				r.Body = ioutil.NopCloser(gzr)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case "deflate":
			r.Body = flate.NewReader(r.Body)
		case "identity":
			// r.Body = r.Body
		default:
			code := http.StatusUnsupportedMediaType
			msg := fmt.Sprintf("Unsupported content encoding %s", encoding)
			http.Error(w, msg, code)
			return
		}
		r.Header.Set("Content-Encoding", "identity")
		next.ServeHTTP(w, r)
	})
})
