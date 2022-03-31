package gurl

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

func newErrorWithCode(message string, code int) *errorWithCode {
	return &errorWithCode{
		message: message,
		code:    code,
	}
}
