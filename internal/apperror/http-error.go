package apperror

import (
	"fmt"
)

type HttpError struct {
	Code    int
	Message string
}

func NewHttpError(code int, message string) *HttpError {
	return &HttpError{
		Code:    code,
		Message: message,
	}
}

func NewBadRequest(msg string, err error) *HttpError {
	return &HttpError{Code: 400, Message: msg}
}

func NewNotFound(msg string, err error) *HttpError {
	return &HttpError{Code: 404, Message: msg}
}

func NewInternal(msg string, err error) *HttpError {
	return &HttpError{Code: 500, Message: msg}
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("status=%d message=%s", e.Code, e.Message)
}
