package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mw "github.com/alxarch/go-http/middleware"
)

func Test_Compose(t *testing.T) {
	ms := make([]mw.Middleware, 0, 4)
	for i := 0; i < 0; i++ {
		ms = append(ms, mw.MiddlewareFunc(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Foo", fmt.Sprintf("MW %d", i))
				next.ServeHTTP(w, r)
			})

		}))

	}
	h := mw.Compose(ms...)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rec, req)
	if hs, ok := rec.HeaderMap["Foo"]; ok {
		if !reflect.DeepEqual(hs, []string{
			"MW 1",
			"MW 2",
			"MW 3",
			"MW 4",
		}) {
			t.Error("Objects are not equal")
		}
	}

}

func Test_RequireMethods(t *testing.T) {
	h := mw.Compose(mw.RequireMethods(http.MethodPost, http.MethodPut))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Error("Wrong code")
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code == http.StatusMethodNotAllowed {
		t.Error("Wrong code")
	}

}
func Test_RequireMimeTypes(t *testing.T) {
	h := mw.Compose(mw.RequireMimeTypes("text/plain", "application/json"))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "foo/bar")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Error("Wrong code")
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "text/plain")
	h.ServeHTTP(rec, req)
	if rec.Code == http.StatusUnsupportedMediaType {
		t.Error("Wrong code")
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "text/plain;foo=bar,foo=baz")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Error("Wrong code")
	}
}
