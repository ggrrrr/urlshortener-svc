package application

import (
	"fmt"
	"net/http"
)

type AppError struct {
	code    int
	message string
	err     error
}

func NewUnauthorized() *AppError {
	return &AppError{
		code:    http.StatusUnauthorized,
		message: "unauthorized request",
	}
}

func NewForbidden(msg string) *AppError {
	return &AppError{
		code:    http.StatusForbidden,
		message: msg,
	}
}

func NewNotFound() *AppError {
	return &AppError{
		code: http.StatusNotFound,
		// message: "",
	}
}

func NewBadRequest(msg string, err error) *AppError {
	return &AppError{
		code:    http.StatusBadRequest,
		message: msg,
		err:     err,
	}
}

func NewSystemError(msg string, err error) *AppError {
	return &AppError{
		code:    http.StatusInternalServerError,
		message: msg,
		err:     err,
	}
}

func (e *AppError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("app.error[%d]: %+v", e.code, e.err)
	}
	return fmt.Sprintf("app.error[%d]: %s: %+v", e.code, e.message, e.err)
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Unwrap() error {
	return e.err
}

func (e *AppError) Message() string {
	return e.message
}
