package httperror

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPError interface {
	Error() string
	Code() int
}

type Error struct {
	code int
	error
}

var _ HTTPError = Error{}

func (e Error) Code() int {
	return e.code
}
func New(args ...interface{}) error {
	code, err := errorFromMsgAndArgs(args)
	return Error{code, err}
}
func Wrap(err error, code int) *Error {
	if err == nil {
		return nil
	}
	return &Error{code, err}
}

func (e *Error) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	Serve(w, e)
}

func defaultErr(code int) error {
	if status := http.StatusText(code); status != "" {
		return errors.New(status)
	}
	return fmt.Errorf("HTTP Status %d", code)
}

func errorFromMsgAndArgs(args []interface{}) (code int, err error) {
	code = http.StatusInternalServerError
	if len(args) == 0 || args == nil {
		return code, defaultErr(code)
	}
	switch e := args[0].(type) {
	case HTTPError:
		return e.Code(), e
	case int:
		code = e
		if len(args) == 1 {
			return code, defaultErr(code)
		} else {
			_, err = errorFromMsgAndArgs(args[1:])
			return code, err
		}
	case error:
		return code, e
	case string:
		if len(args) == 1 {
			return code, errors.New(e)
		} else {
			return code, fmt.Errorf(e, args[1:]...)
		}
	default:
		return code, errors.New(fmt.Sprint(e))
	}
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

func StdHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h.ServeHTTP(w, r)
		Serve(w, err)
	})
}
func Serve(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	if e, isHTTP := err.(HTTPError); isHTTP {
		http.Error(w, e.Error(), e.Code())
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
