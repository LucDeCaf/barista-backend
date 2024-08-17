package errors

import "fmt"

type HttpError struct {
	Code    int
	Message string
}

func NewHttpError(code int, message string) *HttpError {
	return &HttpError{Code: code, Message: message}
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%v %v", e.Code, e.Message)
}
