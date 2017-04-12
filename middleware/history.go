package middleware

import (
	"bytes"
	"container/ring"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"sync"
)

type History struct {
	ring *ring.Ring
	mu   sync.RWMutex
}

type dump struct {
	Request  []byte
	Response []byte
}

func (r *History) Do(do func(req, res []byte)) {
	if nil == r.ring {
		return
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.ring.Do(func(x interface{}) {
		if d, ok := x.(dump); ok {
			do(d.Request, d.Response)
		}
	})
}

func (r *History) Push(req, res []byte) {
	if nil == r.ring {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if d, ok := r.ring.Value.(dump); ok {
		putBuffer(bytes.NewBuffer(d.Request))
		putBuffer(bytes.NewBuffer(d.Response))
	}
	r.ring.Value = dump{req, res}
	r.ring = r.ring.Next()
}

func NewHistory(size int) *History {
	if size <= 0 {
		return &History{
			ring: nil,
		}
	}
	return &History{
		ring: ring.New(size),
	}
}

func (h *History) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req, resp []byte
		defer func() {
			h.Push(req, resp)
		}()
		if dump, err := httputil.DumpRequest(r, true); err == nil {
			req = dump
		}

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		if res := rec.Result(); res != nil {
			if dump, err := httputil.DumpResponse(res, true); err == nil {
				resp = dump
			}
		}

		for name, value := range rec.Header() {
			w.Header()[name] = value
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	})
}
