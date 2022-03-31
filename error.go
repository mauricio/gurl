package gurl

import "fmt"

type ReturnCodeError interface {
	error
	Code() int
}

type errorWithCode struct {
	message string
	code    int
}

func (e *errorWithCode) Error() string {
	return e.message
}

func (e *errorWithCode) Code() int {
	return e.code
}

func newErrorWithCode(code int, message string, args ...any) *errorWithCode {
	return &errorWithCode{
		message: fmt.Sprintf(message, args...),
		code:    code,
	}
}
