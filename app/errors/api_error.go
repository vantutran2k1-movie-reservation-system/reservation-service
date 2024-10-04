package errors

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Error() string {
	return e.Message
}

func BadRequestError(format string, args ...any) *ApiError {
	return newError(http.StatusBadRequest, fmt.Sprintf(format, args...))
}

func UnauthorizedError(format string, args ...any) *ApiError {
	return newError(http.StatusUnauthorized, fmt.Sprintf(format, args...))
}

func NotFoundError(format string, args ...any) *ApiError {
	return newError(http.StatusNotFound, fmt.Sprintf(format, args...))
}

func InternalServerError(format string, args ...any) *ApiError {
	return newError(http.StatusInternalServerError, fmt.Sprintf(format, args...))
}

func newError(statusCode int, message string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Message:    message,
	}
}
