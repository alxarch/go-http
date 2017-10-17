package assert

import (
	"mime"
	"net/http"

	"github.com/alxarch/go-http/httperror"
)

func OK(ok bool, args ...interface{}) {
	if !ok {
		panic(httperror.New(args...))
	}
}

func NoError(err error, code int) {
	if err != nil {
		panic(httperror.New(code, err))
	}
}
func NotNil(x interface{}, args ...interface{}) {
	if x == nil {
		panic(httperror.New(args...))
	}
}

func Method(m string, allow ...string) string {
	for _, a := range allow {
		if m == a {
			return m
		}
	}
	panic(httperror.New(http.StatusMethodNotAllowed))
}

func ContentType(v string, allow ...string) (mediatype string, params map[string]string) {
	var err error
	mediatype, params, err = mime.ParseMediaType(v)
	NoError(err, http.StatusBadRequest)
	for _, a := range allow {
		if mediatype == a {
			return
		}
	}
	panic(httperror.New(http.StatusUnsupportedMediaType))
}
