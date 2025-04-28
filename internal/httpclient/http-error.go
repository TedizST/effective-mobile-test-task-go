package httpclient

import "fmt"

type HttpError struct {
	Method     string
	StatusCode int
	Body       string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("method=%s HTTP statusCode=%d: %s", e.Method, e.StatusCode, e.Body)
}
