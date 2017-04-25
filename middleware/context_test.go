package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/alxarch/go-http/middleware"
)

func Test_Context(t *testing.T) {
	key := struct{ string }{"key"}
	val := struct{ string }{"val"}
	ctx := context.WithValue(context.Background(), key, val)

	cm := middleware.Compose(middleware.Context(ctx), middleware.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !reflect.DeepEqual(r.Context(), ctx) {
				t.Error("Not equal")
			}
			next.ServeHTTP(w, r)
		})
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	cm.ServeHTTP(rec, req)

}
