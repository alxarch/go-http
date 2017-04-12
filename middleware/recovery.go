package middleware

import (
	"net/http"

	"github.com/alxarch/go-http/httperror"
	"github.com/urfave/negroni"
)

func NewRecovery() *Recovery {
	return &Recovery{Recovery: negroni.NewRecovery()}
}

type Recovery struct {
	*negroni.Recovery
}

func (rec *Recovery) Wrap(next http.Handler) http.Handler {
	return Negroni(rec.Recovery).Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				if hp, ok := e.(httperror.HTTPError); ok {
					// Handle http coded errors gracefully
					http.Error(w, hp.Error(), hp.Code())
				} else {
					panic(e)
				}
			}
		}()
		next.ServeHTTP(w, r)
	}))
}
