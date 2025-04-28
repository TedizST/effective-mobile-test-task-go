package apperror

import "fmt"

type AppError struct {
	Method  string
	Message string
	Err     error
}

func NewAppError(method string, message string, err error) *AppError {
	return &AppError{
		Method:  method,
		Message: message,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("method=%s message=%s error=%v", e.Method, e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}
