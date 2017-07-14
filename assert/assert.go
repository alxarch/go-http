package assert

import (
	"mime"
	"net/http"
)

type Error struct {
	StatusCode int
	Message    string
}

func (p Error) Error() string {
	if "" == p.Message {
		return http.StatusText(p.StatusCode)
	}
	return p.Message
}

func (p Error) Code() int {
	return p.StatusCode
}

func OK(ok bool, code int, message string) {
	if !ok {
		panic(Error{code, message})
	}
}

func NoError(err error, code int) {
	if err != nil {
		panic(Error{code, err.Error()})
	}
}

func Panic(code int, message string) {
	panic(Error{code, message})
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Method(m string, allow ...string) string {
	for _, a := range allow {
		if m == a {
			return m
		}
	}
	panic(Error{http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed)})
}

func ContentType(v string, allow ...string) (string, map[string]string) {
	mt, params, err := mime.ParseMediaType(v)
	NoError(err, http.StatusBadRequest)
	for _, a := range allow {
		if mt == a {
			return mt, params
		}
	}
	panic(Error{http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType)})
}
