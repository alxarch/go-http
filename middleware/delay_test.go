package middleware_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alxarch/go-http/middleware"
)

func Test_Delay(t *testing.T) {

	dt := 50 * time.Millisecond
	d := middleware.Delay(dt)
	h := middleware.Compose(d)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	defer func(start time.Time) {
		if time.Now().Sub(start) < dt {
			t.Error("Delay not")
		}

	}(time.Now())
	h.ServeHTTP(rec, req)

}
