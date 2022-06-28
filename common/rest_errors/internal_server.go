package resterrors

import "strconv"

type InternalServerError struct {
	StatusCode int
	ErrMsg     string
}

func NewInternalServerError(errMsg string) *InternalServerError {
	return &InternalServerError{StatusCode: 500, ErrMsg: errMsg}
}

func (e *InternalServerError) HttpStatusCode() int {
	return e.StatusCode
}

func (e *InternalServerError) HttpStatusCodeString() string {
	return strconv.Itoa(e.StatusCode)
}

func (e *InternalServerError) Error() string {
	return e.ErrMsg
}
