package httperror

type HTTPError interface {
	Error() string
	Code() int
}
