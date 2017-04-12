package assert

import "net/http"

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
