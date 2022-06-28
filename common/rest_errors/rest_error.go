package resterrors

type RestError interface {
	HttpStatusCode() int
	HttpStatusCodeString() string
	Error() string
}
