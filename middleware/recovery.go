package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/alxarch/go-http/httperror"
)

func NewRecovery() *Recovery {
	return &Recovery{
		Logger:    log.New(os.Stdout, "", 0),
		StackSize: DefaultBufferSize,
	}
}

type Recovery struct {
	Logger     *log.Logger
	PrintStack bool
	StackAll   bool
	StackSize  int
}

func (rec *Recovery) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				if hp, ok := e.(httperror.HTTPError); ok {
					// Handle http coded errors gracefully
					http.Error(w, hp.Error(), hp.Code())
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					f := "PANIC: %s\n%s"
					var stack []byte
					if rec.StackSize > 0 {
						stack = make([]byte, rec.StackSize)
						stack = stack[:runtime.Stack(stack, rec.StackAll)]
					} else {
						stack = []byte{}
					}
					if rec.PrintStack {
						fmt.Fprintf(w, f, e, stack)
					}
					if rec.Logger != nil {
						rec.Logger.Printf(f, e, stack)
					}
				}
			}
		}()
		next.ServeHTTP(w, r)
	})

}
