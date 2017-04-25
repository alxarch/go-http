package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alxarch/go-http/middleware"
)

func Test_Timeout(t *testing.T) {
	tmw := middleware.Timeout(50*time.Millisecond, nil)
	delay := middleware.Delay(10 * time.Millisecond)
	h := middleware.Compose(tmw, delay)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Wrong code %d", rec.Code)
	}

	delay = middleware.Delay(60 * time.Millisecond)
	h = middleware.Compose(tmw, delay)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusGatewayTimeout {
		t.Errorf("Wrong code %d", rec.Code)
	}
}
