package resterrors

import "strconv"

type BadReqeustError struct {
	StatusCode int
	ErrMsg     string
}

func NewBadRequestError(errMsg string) *BadReqeustError {
	return &BadReqeustError{StatusCode: 400, ErrMsg: errMsg}
}

func (e *BadReqeustError) HttpStatusCode() int {
	return e.StatusCode
}

func (e *BadReqeustError) HttpStatusCodeString() string {
	return strconv.Itoa(e.StatusCode)
}

func (e *BadReqeustError) Error() string {
	return e.ErrMsg
}
