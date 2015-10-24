package controller

type HTTPError interface {
	error
	Code() int
}

type httpError struct {
	code int
	msg string
}

func NewHTTPError(code int, msg string) httpError {
	return httpError{
		code: code,
		msg: msg,
	}
}

func (e httpError) Error() string {
	return e.msg
}

func (e httpError) Code() int {
	return e.code
}
