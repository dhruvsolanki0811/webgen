package domain

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func ErrBadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func ErrUnauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func ErrForbidden(msg string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: msg}
}

func ErrNotFound(resource string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: fmt.Sprintf("%s not found", resource)}
}

func ErrConflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}

func ErrTooManyRequests() *AppError {
	return &AppError{Code: http.StatusTooManyRequests, Message: "rate limit exceeded"}
}

func ErrInternal(err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: "internal error", Err: err}
}
