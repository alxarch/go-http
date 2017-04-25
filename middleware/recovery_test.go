package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alxarch/go-http/assert"
	"github.com/alxarch/go-http/middleware"
)

func Test_RecoveryHTTPError(t *testing.T) {
	rec := middleware.NewRecovery()
	h := rec.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Panic(http.StatusExpectationFailed, "Foo bar.")
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, r)
	if w.Code != http.StatusExpectationFailed {
		t.Error("Wrong code %d", w.Code)
	}

}
func Test_Recovery(t *testing.T) {
	rec := middleware.NewRecovery()
	rec.PrintStack = true
	h := rec.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Foo bar.")
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Error("Wrong code %d", w.Code)
	}

}
